package seeder

import (
	cloudflarer2 "POS-kasir/pkg/cloudflare-r2"
	"POS-kasir/pkg/logger"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedProductionData(ctx context.Context, pool *pgxpool.Pool, r2Client cloudflarer2.IR2, log logger.ILogger) error {
	log.Info("Seeding production data...")

	// Check if production data already exists (idempotent)
	var orderCount int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM orders").Scan(&orderCount)
	if err != nil {
		return fmt.Errorf("failed to check existing orders: %w", err)
	}
	if orderCount > 100 {
		log.Info("Production data already exists (found >100 orders), skipping.")
		return nil
	}

	// Get user IDs for assigning orders
	userIDs, err := getUserIDs(ctx, pool)
	if err != nil {
		return fmt.Errorf("failed to get user IDs: %w", err)
	}
	if len(userIDs) == 0 {
		return fmt.Errorf("no users found, please seed users first")
	}

	// Get category IDs
	categoryMap, err := getCategoryMap(ctx, pool)
	if err != nil {
		return fmt.Errorf("failed to get categories: %w", err)
	}

	// Get payment method IDs
	paymentMethodMap, err := getPaymentMethodMap(ctx, pool)
	if err != nil {
		return fmt.Errorf("failed to get payment methods: %w", err)
	}

	// Get cancellation reasons
	cancelReasons, err := getCancellationReasons(ctx, pool)
	if err != nil {
		return fmt.Errorf("failed to get cancellation reasons: %w", err)
	}

	// Get promotions
	promotions, err := getPromotions(ctx, pool)
	if err != nil {
		return fmt.Errorf("failed to get promotions: %w", err)
	}

	// 1. Seed Products + Options + Category Assignments
	products, err := seedProducts(ctx, pool, categoryMap, r2Client, log)
	if err != nil {
		return fmt.Errorf("failed to seed products: %w", err)
	}
	log.Infof("Seeded %d products with variants", len(products))

	// 2. Seed Customers
	customerIDs, err := seedCustomers(ctx, pool, log)
	if err != nil {
		return fmt.Errorf("failed to seed customers: %w", err)
	}
	log.Infof("Seeded %d customers", len(customerIDs))

	// 3. Seed Promotions
	err = seedPromotions(ctx, pool, log)
	if err != nil {
		return fmt.Errorf("failed to seed promotions: %w", err)
	}
	log.Info("Seeded promotions")

	// 4. Seed Orders (with backdated timestamps)
	totalOrders, err := seedOrders(ctx, pool, products, userIDs, customerIDs, paymentMethodMap, cancelReasons, promotions, log)
	if err != nil {
		return fmt.Errorf("failed to seed orders: %w", err)
	}
	log.Infof("Seeded %d orders spanning 2-4 months", totalOrders)

	// 5. Seed Shifts
	if err := seedShifts(ctx, pool, userIDs, log); err != nil {
		return fmt.Errorf("failed to seed shifts: %w", err)
	}
	log.Info("Seeded cashier shifts")

	log.Info("Production data seeding completed successfully!")
	return nil
}

// ─── Helper Types ────────────────────────────────────────────────

func uploadImageFromUrl(ctx context.Context, r2Client cloudflarer2.IR2, urlStr string, folder string, id string, log logger.ILogger) *string {
	if urlStr == "" || r2Client == nil {
		return &urlStr
	}

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		log.Warnf("Failed to create request for %s: %v", urlStr, err)
		return &urlStr
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Warnf("Failed to download image from %s: %v", urlStr, err)
		return &urlStr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warnf("Failed to download image from %s, status: %d", urlStr, resp.StatusCode)
		return &urlStr
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("Failed to read image data from %s: %v", urlStr, err)
		return &urlStr
	}

	objectName := fmt.Sprintf("%s/%s.jpg", folder, id)
	_, err = r2Client.UploadFile(ctx, objectName, data, "image/jpeg")
	if err != nil {
		log.Warnf("Failed to upload image to object storage for %s: %v", objectName, err)
		return &urlStr
	}

	// Returning the object key (e.g. products/UUID.jpg) instead of the full URL,
	// because that's what the application expects to store in the DB.
	return &objectName
}

type seededProduct struct {
	ID        string
	Name      string
	Price     int64
	CostPrice int64
	Options   []seededOption
}

type seededOption struct {
	ID              string
	Name            string
	AdditionalPrice int64
}

// ─── Helper Functions ────────────────────────────────────────────

func getUserIDs(ctx context.Context, pool *pgxpool.Pool) ([]string, error) {
	rows, err := pool.Query(ctx, "SELECT id FROM users WHERE is_active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func getCategoryMap(ctx context.Context, pool *pgxpool.Pool) (map[string]int32, error) {
	rows, err := pool.Query(ctx, "SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := make(map[string]int32)
	for rows.Next() {
		var id int32
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		m[name] = id
	}
	return m, rows.Err()
}

func getPaymentMethodMap(ctx context.Context, pool *pgxpool.Pool) (map[string]int32, error) {
	rows, err := pool.Query(ctx, "SELECT id, name FROM payment_methods WHERE is_active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := make(map[string]int32)
	for rows.Next() {
		var id int32
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		m[name] = id
	}
	return m, rows.Err()
}

func getCancellationReasons(ctx context.Context, pool *pgxpool.Pool) ([]int32, error) {
	rows, err := pool.Query(ctx, "SELECT id FROM cancellation_reasons")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func getPromotions(ctx context.Context, pool *pgxpool.Pool) ([]string, error) {
	rows, err := pool.Query(ctx, "SELECT id FROM promotions WHERE is_active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// ─── 1. Products ─────────────────────────────────────────────────

func seedProducts(ctx context.Context, pool *pgxpool.Pool, categoryMap map[string]int32, r2Client cloudflarer2.IR2, log logger.ILogger) ([]seededProduct, error) {
	type productDef struct {
		Name      string
		Price     int64
		CostPrice int64
		Stock     int32
		ImageURL  string
		Category  string
		Options   []struct {
			Name            string
			AdditionalPrice int64
			ImageURL        string
		}
	}

	defs := []productDef{
		// ─── Makanan ───
		{Name: "Nasi Goreng Spesial", Price: 25000, CostPrice: 12000, Stock: 200, Category: "Makanan",
			ImageURL: "https://images.unsplash.com/photo-1512058564366-18510be2db19?w=400&h=400&fit=crop",
			Options: []struct {
				Name            string
				AdditionalPrice int64
				ImageURL        string
			}{
				{Name: "Pedas", AdditionalPrice: 0, ImageURL: "https://images.unsplash.com/photo-1512058564366-18510be2db19?w=200&h=200&fit=crop"},
				{Name: "Ekstra Telur", AdditionalPrice: 5000, ImageURL: "https://images.unsplash.com/photo-1482049016688-2d3e1b311543?w=200&h=200&fit=crop"},
				{Name: "Ekstra Ayam", AdditionalPrice: 8000, ImageURL: "https://images.unsplash.com/photo-1598515214211-89d3c73ae83b?w=200&h=200&fit=crop"},
			}},
		{Name: "Mie Goreng", Price: 22000, CostPrice: 10000, Stock: 180, Category: "Makanan",
			ImageURL: "https://images.unsplash.com/photo-1585032226651-759b368d7246?w=400&h=400&fit=crop",
			Options: []struct {
				Name            string
				AdditionalPrice int64
				ImageURL        string
			}{
				{Name: "Pedas", AdditionalPrice: 0, ImageURL: "https://images.unsplash.com/photo-1585032226651-759b368d7246?w=200&h=200&fit=crop"},
				{Name: "Ekstra Bakso", AdditionalPrice: 5000, ImageURL: "https://images.unsplash.com/photo-1529692236671-f1f6cf9683ba?w=200&h=200&fit=crop"},
			}},
		{Name: "Mie Ayam Bakso", Price: 20000, CostPrice: 9000, Stock: 150, Category: "Makanan",
			ImageURL: "https://images.unsplash.com/photo-1569718212165-3a8278d5f624?w=400&h=400&fit=crop"},
		{Name: "Ayam Goreng Kremes", Price: 28000, CostPrice: 14000, Stock: 120, Category: "Makanan",
			ImageURL: "https://images.unsplash.com/photo-1626645738196-c2a7c87a8f58?w=400&h=400&fit=crop",
			Options: []struct {
				Name            string
				AdditionalPrice int64
				ImageURL        string
			}{
				{Name: "Sambal Matah", AdditionalPrice: 3000, ImageURL: "https://images.unsplash.com/photo-1563379926898-05f4575a45d8?w=200&h=200&fit=crop"},
				{Name: "Lalapan", AdditionalPrice: 2000, ImageURL: "https://images.unsplash.com/photo-1540420773420-3366772f4999?w=200&h=200&fit=crop"},
			}},
		{Name: "Soto Ayam", Price: 18000, CostPrice: 8000, Stock: 160, Category: "Makanan",
			ImageURL: "https://images.unsplash.com/photo-1547592166-23ac45744acd?w=400&h=400&fit=crop"},
		{Name: "Nasi Rendang", Price: 32000, CostPrice: 16000, Stock: 100, Category: "Makanan",
			ImageURL: "https://images.unsplash.com/photo-1604908176997-125f25cc6f3d?w=400&h=400&fit=crop"},
		{Name: "Gado-Gado", Price: 15000, CostPrice: 7000, Stock: 3, Category: "Makanan", // Intentionally low stock
			ImageURL: "https://images.unsplash.com/photo-1512058564366-18510be2db19?w=400&h=400&fit=crop&q=60"},
		{Name: "Bakso Urat", Price: 20000, CostPrice: 9500, Stock: 140, Category: "Makanan",
			ImageURL: "https://images.unsplash.com/photo-1529692236671-f1f6cf9683ba?w=400&h=400&fit=crop",
			Options: []struct {
				Name            string
				AdditionalPrice int64
				ImageURL        string
			}{
				{Name: "Jumbo", AdditionalPrice: 5000, ImageURL: "https://images.unsplash.com/photo-1529692236671-f1f6cf9683ba?w=200&h=200&fit=crop"},
			}},
		{Name: "Nasi Uduk Komplit", Price: 22000, CostPrice: 10000, Stock: 130, Category: "Makanan",
			ImageURL: "https://images.unsplash.com/photo-1596797038530-2c107229654b?w=400&h=400&fit=crop"},
		{Name: "Sate Ayam (10 tusuk)", Price: 25000, CostPrice: 12000, Stock: 110, Category: "Makanan",
			ImageURL: "https://images.unsplash.com/photo-1555939594-58d7cb561ad1?w=400&h=400&fit=crop"},

		// ─── Minuman ───
		{Name: "Es Teh Manis", Price: 5000, CostPrice: 1500, Stock: 500, Category: "Minuman",
			ImageURL: "https://images.unsplash.com/photo-1556679343-c7306c1976bc?w=400&h=400&fit=crop",
			Options: []struct {
				Name            string
				AdditionalPrice int64
				ImageURL        string
			}{
				{Name: "Large", AdditionalPrice: 3000, ImageURL: "https://images.unsplash.com/photo-1556679343-c7306c1976bc?w=200&h=200&fit=crop"},
			}},
		{Name: "Es Jeruk", Price: 7000, CostPrice: 2500, Stock: 400, Category: "Minuman",
			ImageURL: "https://images.unsplash.com/photo-1621263764928-df1444c5e859?w=400&h=400&fit=crop"},
		{Name: "Kopi Susu Gula Aren", Price: 18000, CostPrice: 6000, Stock: 300, Category: "Minuman",
			ImageURL: "https://images.unsplash.com/photo-1461023058943-07fcbe16d735?w=400&h=400&fit=crop",
			Options: []struct {
				Name            string
				AdditionalPrice int64
				ImageURL        string
			}{
				{Name: "Iced", AdditionalPrice: 0, ImageURL: "https://images.unsplash.com/photo-1461023058943-07fcbe16d735?w=200&h=200&fit=crop"},
				{Name: "Hot", AdditionalPrice: 0, ImageURL: "https://images.unsplash.com/photo-1495774856032-8b90bbb32b32?w=200&h=200&fit=crop"},
				{Name: "Extra Shot", AdditionalPrice: 5000, ImageURL: "https://images.unsplash.com/photo-1509042239860-f550ce710b93?w=200&h=200&fit=crop"},
			}},
		{Name: "Jus Alpukat", Price: 15000, CostPrice: 5000, Stock: 250, Category: "Minuman",
			ImageURL: "https://images.unsplash.com/photo-1623065422902-30a2d299bbe4?w=400&h=400&fit=crop"},
		{Name: "Air Mineral", Price: 4000, CostPrice: 1000, Stock: 600, Category: "Minuman",
			ImageURL: "https://images.unsplash.com/photo-1548839140-29a749e1cf4d?w=400&h=400&fit=crop"},

		// ─── Camilan ───
		{Name: "Kentang Goreng", Price: 15000, CostPrice: 6000, Stock: 200, Category: "Camilan",
			ImageURL: "https://images.unsplash.com/photo-1573080496219-bb080dd4f877?w=400&h=400&fit=crop",
			Options: []struct {
				Name            string
				AdditionalPrice int64
				ImageURL        string
			}{
				{Name: "Cheese", AdditionalPrice: 5000, ImageURL: "https://images.unsplash.com/photo-1573080496219-bb080dd4f877?w=200&h=200&fit=crop"},
			}},
		{Name: "Pisang Goreng Keju", Price: 12000, CostPrice: 5000, Stock: 180, Category: "Camilan",
			ImageURL: "https://images.unsplash.com/photo-1528735602780-2552fd46c7af?w=400&h=400&fit=crop"},
		{Name: "Tahu Crispy", Price: 10000, CostPrice: 3500, Stock: 8, Category: "Camilan", // Intentionally low stock
			ImageURL: "https://images.unsplash.com/photo-1585032226651-759b368d7246?w=400&h=400&fit=crop&q=70"},

		// ─── Makanan Penutup ───
		{Name: "Es Campur", Price: 12000, CostPrice: 4000, Stock: 5, Category: "Makanan Penutup", // Intentionally low stock
			ImageURL: "https://images.unsplash.com/photo-1551024506-0bccd828d307?w=400&h=400&fit=crop"},
		{Name: "Kolak Pisang", Price: 10000, CostPrice: 3500, Stock: 150, Category: "Makanan Penutup",
			ImageURL: "https://images.unsplash.com/photo-1528735602780-2552fd46c7af?w=400&h=400&fit=crop"},

		// ─── Paket ───
		{Name: "Paket Nasi Goreng + Es Teh", Price: 27000, CostPrice: 13000, Stock: 100, Category: "Paket",
			ImageURL: "https://images.unsplash.com/photo-1512058564366-18510be2db19?w=400&h=400&fit=crop&q=80"},
		{Name: "Paket Ayam Goreng + Jus", Price: 38000, CostPrice: 18000, Stock: 80, Category: "Paket",
			ImageURL: "https://images.unsplash.com/photo-1626645738196-c2a7c87a8f58?w=400&h=400&fit=crop&q=80"},
	}

	var products []seededProduct

	for _, d := range defs {
		var productID string
		err := pool.QueryRow(ctx,
			`INSERT INTO products (name, image_url, price, stock, cost_price)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT DO NOTHING
			 RETURNING id`,
			d.Name, d.ImageURL, d.Price, d.Stock, d.CostPrice,
		).Scan(&productID)
		if err != nil {
			log.Infof("Product '%s' may already exist, skipping: %v", d.Name, err)
			continue
		}

		// Download and upload image if URL exists
		storedImageURL := uploadImageFromUrl(ctx, r2Client, d.ImageURL, "products", productID, log)
		if storedImageURL != nil && *storedImageURL != d.ImageURL {
			_, err = pool.Exec(ctx, "UPDATE products SET image_url = $1 WHERE id = $2", *storedImageURL, productID)
			if err != nil {
				log.Warnf("Failed to update product %s with storage image URL: %v", productID, err)
			}
		}

		// Assign category
		if catID, ok := categoryMap[d.Category]; ok {
			_, _ = pool.Exec(ctx,
				`INSERT INTO product_categories (product_id, category_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
				productID, catID,
			)
		}

		p := seededProduct{
			ID:        productID,
			Name:      d.Name,
			Price:     d.Price,
			CostPrice: d.CostPrice,
		}

		// Seed options
		for _, opt := range d.Options {
			var optID string
			err := pool.QueryRow(ctx,
				`INSERT INTO product_options (product_id, name, additional_price, image_url)
				 VALUES ($1, $2, $3, $4)
				 RETURNING id`,
				productID, opt.Name, opt.AdditionalPrice, opt.ImageURL,
			).Scan(&optID)
			if err != nil {
				log.Infof("Option '%s' for product '%s' failed: %v", opt.Name, d.Name, err)
				continue
			}

			// Download and upload option image if URL exists
			storedOptImageURL := uploadImageFromUrl(ctx, r2Client, opt.ImageURL, "product_options", optID, log)
			if storedOptImageURL != nil && *storedOptImageURL != opt.ImageURL {
				_, err = pool.Exec(ctx, "UPDATE product_options SET image_url = $1 WHERE id = $2", *storedOptImageURL, optID)
				if err != nil {
					log.Warnf("Failed to update product option %s with storage image URL: %v", optID, err)
				}
			}

			p.Options = append(p.Options, seededOption{
				ID:              optID,
				Name:            opt.Name,
				AdditionalPrice: opt.AdditionalPrice,
			})
		}

		products = append(products, p)
	}

	return products, nil
}

// ─── 2. Customers ────────────────────────────────────────────────

func seedCustomers(ctx context.Context, pool *pgxpool.Pool, log logger.ILogger) ([]string, error) {
	customers := []struct {
		Name    string
		Phone   string
		Email   string
		Address string
	}{
		{"Budi Santoso", "081234567890", "budi.santoso@email.com", "Jl. Sudirman No. 10, Jakarta"},
		{"Siti Rahayu", "081234567891", "siti.rahayu@email.com", "Jl. Gatot Subroto No. 5, Jakarta"},
		{"Ahmad Hidayat", "081234567892", "ahmad.hidayat@email.com", "Jl. MH Thamrin No. 15, Jakarta"},
		{"Dewi Lestari", "081234567893", "dewi.lestari@email.com", "Jl. Kebon Sirih No. 8, Jakarta"},
		{"Rizky Pratama", "081234567894", "rizky.pratama@email.com", "Jl. Cikini Raya No. 22, Jakarta"},
		{"Nur Aini", "081234567895", "nur.aini@email.com", "Jl. Veteran No. 3, Bandung"},
		{"Hendra Wijaya", "081234567896", "hendra.wijaya@email.com", "Jl. Asia Afrika No. 12, Bandung"},
		{"Maya Sari", "081234567897", "maya.sari@email.com", "Jl. Malioboro No. 1, Yogyakarta"},
		{"Fajar Nugroho", "081234567898", "fajar.nugroho@email.com", "Jl. Diponegoro No. 7, Surabaya"},
		{"Putri Amelia", "081234567899", "putri.amelia@email.com", "Jl. Pemuda No. 20, Semarang"},
	}

	var ids []string
	for _, c := range customers {
		var id string
		err := pool.QueryRow(ctx,
			`INSERT INTO customers (name, phone, email, address)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT DO NOTHING
			 RETURNING id`,
			c.Name, c.Phone, c.Email, c.Address,
		).Scan(&id)
		if err != nil {
			log.Infof("Customer '%s' may already exist: %v", c.Name, err)
			continue
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// ─── 3. Promotions ───────────────────────────────────────────────

func seedPromotions(ctx context.Context, pool *pgxpool.Pool, log logger.ILogger) error {
	now := time.Now()
	threeMonthsAgo := now.AddDate(0, -3, 0)
	oneMonthLater := now.AddDate(0, 1, 0)

	promos := []struct {
		Name          string
		Description   string
		Scope         string
		DiscountType  string
		DiscountValue float64
		MaxDiscount   float64
		StartDate     time.Time
		EndDate       time.Time
		RuleType      string
		RuleValue     string
	}{
		{
			Name: "Diskon Grand Opening", Description: "Diskon 10% untuk semua pesanan di atas Rp50.000",
			Scope: "ORDER", DiscountType: "percentage", DiscountValue: 10, MaxDiscount: 15000,
			StartDate: threeMonthsAgo, EndDate: oneMonthLater,
			RuleType: "MINIMUM_ORDER_AMOUNT", RuleValue: "50000",
		},
		{
			Name: "Promo Paket Hemat", Description: "Potongan Rp5.000 untuk pesanan minimal Rp30.000",
			Scope: "ORDER", DiscountType: "fixed_amount", DiscountValue: 5000, MaxDiscount: 5000,
			StartDate: threeMonthsAgo, EndDate: oneMonthLater,
			RuleType: "MINIMUM_ORDER_AMOUNT", RuleValue: "30000",
		},
		{
			Name: "Happy Hour 15%", Description: "Diskon 15% untuk semua pesanan, maks Rp20.000",
			Scope: "ORDER", DiscountType: "percentage", DiscountValue: 15, MaxDiscount: 20000,
			StartDate: threeMonthsAgo, EndDate: oneMonthLater,
			RuleType: "MINIMUM_ORDER_AMOUNT", RuleValue: "25000",
		},
	}

	for _, p := range promos {
		var promoID string
		err := pool.QueryRow(ctx,
			`INSERT INTO promotions (name, description, scope, discount_type, discount_value, max_discount_amount, start_date, end_date, is_active)
			 VALUES ($1, $2, $3::promotion_scope, $4::discount_type, $5, $6, $7, $8, true)
			 ON CONFLICT DO NOTHING
			 RETURNING id`,
			p.Name, p.Description, p.Scope, p.DiscountType, p.DiscountValue, p.MaxDiscount, p.StartDate, p.EndDate,
		).Scan(&promoID)
		if err != nil {
			log.Infof("Promotion '%s' may already exist: %v", p.Name, err)
			continue
		}

		// Add rule
		_, err = pool.Exec(ctx,
			`INSERT INTO promotion_rules (promotion_id, rule_type, rule_value)
			 VALUES ($1, $2::promotion_rule_type, $3)`,
			promoID, p.RuleType, p.RuleValue,
		)
		if err != nil {
			log.Errorf("Failed to add rule for promotion '%s': %v", p.Name, err)
		}
	}

	return nil
}

// ─── 4. Orders ───────────────────────────────────────────────────

func seedOrders(ctx context.Context, pool *pgxpool.Pool, products []seededProduct, userIDs []string, customerIDs []string, paymentMethodMap map[string]int32, cancelReasons []int32, promotions []string, log logger.ILogger) (int, error) {
	if len(products) == 0 {
		return 0, fmt.Errorf("no products to create orders from")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	now := time.Now()

	// Go back 2-4 months (60-120 days)
	startDaysBack := 60 + rng.Intn(61) // 60-120
	startDate := now.AddDate(0, 0, -startDaysBack)

	orderTypes := []string{"dine_in", "takeaway"}
	paymentMethods := []string{"Cash", "QRIS Dinamis", "QRIS Statis"}
	paymentWeights := []int{70, 20, 10} // Weighted distribution

	totalOrders := 0

	// Iterate each day from startDate to now
	for d := startDate; d.Before(now); d = d.AddDate(0, 0, 1) {
		weekday := d.Weekday()

		// More orders on weekends
		var ordersToday int
		if weekday == time.Saturday || weekday == time.Sunday {
			ordersToday = 5 + rng.Intn(8) // 5-12
		} else {
			ordersToday = 3 + rng.Intn(5) // 3-7
		}

		for i := 0; i < ordersToday; i++ {
			// Random time between 08:00 and 21:00
			hour := 8 + rng.Intn(13)
			minute := rng.Intn(60)
			orderTime := time.Date(d.Year(), d.Month(), d.Day(), hour, minute, rng.Intn(60), 0, d.Location())

			// Random user
			userID := userIDs[rng.Intn(len(userIDs))]

			// Random order type (60% dine_in, 40% takeaway)
			orderType := orderTypes[0]
			if rng.Intn(100) < 40 {
				orderType = orderTypes[1]
			}

			// Random customer (30% chance)
			var customerID *string
			if len(customerIDs) > 0 && rng.Intn(100) < 30 {
				cid := customerIDs[rng.Intn(len(customerIDs))]
				customerID = &cid
			}

			// Pick payment method (weighted)
			pmName := pickWeighted(rng, paymentMethods, paymentWeights)
			pmID, ok := paymentMethodMap[pmName]
			if !ok {
				// fallback to Cash
				for name, id := range paymentMethodMap {
					pmID = id
					pmName = name
					break
				}
			}

			// 1-4 random items
			numItems := 1 + rng.Intn(4)
			usedProducts := make(map[int]bool)
			var grossTotal int64

			type orderItemData struct {
				productID string
				qty       int32
				price     int64
				subtotal  int64
				costPrice int64
			}
			var items []orderItemData

			for j := 0; j < numItems; j++ {
				// Pick a random product (avoid duplicates)
				idx := rng.Intn(len(products))
				if usedProducts[idx] {
					continue
				}
				usedProducts[idx] = true

				product := products[idx]
				qty := int32(1 + rng.Intn(3)) // 1-3

				// Sometimes include an option price
				itemPrice := product.Price
				if len(product.Options) > 0 && rng.Intn(100) < 40 {
					opt := product.Options[rng.Intn(len(product.Options))]
					itemPrice += opt.AdditionalPrice
				}

				subtotal := itemPrice * int64(qty)
				grossTotal += subtotal

				items = append(items, orderItemData{
					productID: product.ID,
					qty:       qty,
					price:     itemPrice,
					subtotal:  subtotal,
					costPrice: product.CostPrice,
				})
			}

			if len(items) == 0 {
				continue
			}

			// ~10% cancelled
			isCancelled := rng.Intn(100) < 10
			status := "paid"
			var cancelReasonID *int32
			var cancelNote *string

			if isCancelled {
				status = "cancelled"
				if len(cancelReasons) > 0 {
					id := cancelReasons[rng.Intn(len(cancelReasons))]
					cancelReasonID = &id
					note := "Dibatalkan secara acak dari seeder"
					cancelNote = &note
				}
			}

			// Apply promotion? (~20% chance, only if not cancelled)
			var appliedPromotionID *string
			var discountAmount int64 = 0
			if !isCancelled && len(promotions) > 0 && rng.Intn(100) < 20 {
				id := promotions[rng.Intn(len(promotions))]
				appliedPromotionID = &id
				// Random discount: either flat 5000 or 10%
				if rng.Intn(2) == 0 {
					discountAmount = 5000
				} else {
					discountAmount = grossTotal * 10 / 100
				}
				if discountAmount > grossTotal {
					discountAmount = grossTotal
				}
			}

			netTotal := grossTotal - discountAmount

			// Cash payment: calculate change
			var cashReceived *int64
			var changeDue *int64
			if pmName == "Cash" && !isCancelled {
				// Round up to nearest 5000 or 10000
				rounded := ((netTotal / 5000) + 1) * 5000
				if rounded < netTotal+1000 {
					rounded += 5000
				}
				cashReceived = &rounded
				change := rounded - netTotal
				changeDue = &change
			}

			// Insert order with backdated timestamp
			var orderID string
			err := pool.QueryRow(ctx,
				`INSERT INTO orders (user_id, type, status, gross_total, discount_amount, net_total,
					payment_method_id, cash_received, change_due, customer_id, 
					cancellation_reason_id, cancellation_notes, applied_promotion_id, 
					created_at, updated_at, version)
				 VALUES ($1, $2::order_type, $3::order_status, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $14, 1)
				 RETURNING id`,
				userID, orderType, status, grossTotal, discountAmount, netTotal, pmID,
				cashReceived, changeDue, customerID, cancelReasonID, cancelNote, appliedPromotionID, orderTime,
			).Scan(&orderID)
			if err != nil {
				log.Errorf("Failed to create order: %v", err)
				continue
			}

			// Insert order items
			for _, item := range items {
				_, err := pool.Exec(ctx,
					`INSERT INTO order_items (order_id, product_id, quantity, price_at_sale, subtotal, net_subtotal, cost_price_at_sale)
					 VALUES ($1, $2, $3, $4, $5, $5, $6)`,
					orderID, item.productID, item.qty, item.price, item.subtotal, item.costPrice,
				)
				if err != nil {
					log.Errorf("Failed to create order item: %v", err)
				}
			}

			totalOrders++
		}
	}

	return totalOrders, nil
}

func pickWeighted(rng *rand.Rand, items []string, weights []int) string {
	total := 0
	for _, w := range weights {
		total += w
	}
	r := rng.Intn(total)
	cumulative := 0
	for i, w := range weights {
		cumulative += w
		if r < cumulative {
			return items[i]
		}
	}
	return items[len(items)-1]
}

// ─── 5. Shifts ───────────────────────────────────────────────────

func seedShifts(ctx context.Context, pool *pgxpool.Pool, userIDs []string, log logger.ILogger) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	now := time.Now()
	startDate := now.AddDate(0, 0, -120) // go back 120 days

	for d := startDate; d.Before(now); d = d.AddDate(0, 0, 1) {
		// Create 1 or 2 shifts per day
		numShifts := 1
		if rng.Intn(100) < 50 {
			numShifts = 2
		}

		for i := 0; i < numShifts; i++ {
			userID := userIDs[rng.Intn(len(userIDs))]
			
			// Shift start time between 07:00 and 14:00
			startHour := 7 + rng.Intn(8)
			startTime := time.Date(d.Year(), d.Month(), d.Day(), startHour, rng.Intn(60), 0, 0, d.Location())
			
			// Shift end time 6-10 hours later
			durationHours := 6 + rng.Intn(5)
			endTime := startTime.Add(time.Duration(durationHours) * time.Hour)
			if endTime.After(now) {
				continue
			}

			startCash := int64(rng.Intn(5)+5) * 100000 // 500k to 1m
			
			// Simulate sales cash (just random for seed purposes, 1m to 5m)
			salesCash := int64(rng.Intn(40)+10) * 100000
			expectedCash := startCash + salesCash
			
			// Actual cash is mostly accurate, sometimes slightly off (+/- 20k)
			actualCash := expectedCash
			if rng.Intn(100) < 20 {
				diff := int64(rng.Intn(40)-20) * 1000 // -20k to +20k
				actualCash += diff
			}

			_, err := pool.Exec(ctx,
				`INSERT INTO shifts (user_id, start_time, end_time, start_cash, expected_cash_end, actual_cash_end, status, created_at, updated_at)
				 VALUES ($1, $2, $3, $4, $5, $6, 'closed', $2, $3)`,
				userID, startTime, endTime, startCash, expectedCash, actualCash,
			)
			if err != nil {
				log.Errorf("Failed to seed shift: %v", err)
			}
		}
	}
	return nil
}
