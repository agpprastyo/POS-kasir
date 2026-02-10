
(function () {
    const POLL_INTERVAL = 500;
    const MAX_RETRIES = 20;
    let attempts = 0;

    const initRBAC = () => {
        // Wait until UI system is ready
        if (typeof window.ui === 'undefined' || typeof window.ui.spec !== 'function') {
            if (attempts < MAX_RETRIES) {
                attempts++;
                setTimeout(initRBAC, POLL_INTERVAL);
            }
            return;
        }

        // Observer to handle UI updates (e.g. expanding tags, filtering)
        const observer = new MutationObserver((mutations) => {
            enrichUI();
        });

        const app = document.getElementById('swagger-ui');
        if (app) {
            observer.observe(app, {
                childList: true,
                subtree: true
            });
            // Initial attempt
            setTimeout(enrichUI, 1000);
        }
    };

    const enrichUI = () => {
        try {
            // Get spec from Swagger UI system
            // Try compatible selector access
            let spec = null;
            if (window.ui.spec().selectors && window.ui.spec().selectors.specJson) {
                spec = window.ui.spec().selectors.specJson().toJS();
            } else if (window.ui.getSystem().specSelectors && window.ui.getSystem().specSelectors.specJson) {
                spec = window.ui.getSystem().specSelectors.specJson().toJS();
            }

            if (!spec || !spec.paths) return;

            const opBlocks = document.querySelectorAll('.opblock');

            opBlocks.forEach(block => {
                if (block.dataset.rbacEnriched === 'true') return;

                const methodEl = block.querySelector('.opblock-summary-method');
                // Use data-path attribute for reliable path matching
                const pathEl = block.querySelector('.opblock-summary-path');

                if (!methodEl || !pathEl) return;

                const method = methodEl.textContent.trim().toLowerCase();
                const path = pathEl.getAttribute('data-path');

                if (!path) return;

                // Look up in spec
                if (spec.paths[path] && spec.paths[path][method]) {
                    const op = spec.paths[path][method];
                    if (op['x-roles'] && Array.isArray(op['x-roles'])) {
                        createBadge(block, op['x-roles']);
                    }
                }

                block.dataset.rbacEnriched = 'true';
            });
        } catch (e) {
            console.warn('RBAC Enrichment error:', e);
        }
    };

    const createBadge = (block, roles) => {
        const summary = block.querySelector('.opblock-summary-operation-id') || block.querySelector('.opblock-summary-description') || block.querySelector('.opblock-summary');

        // Prevent duplicate badges
        if (block.querySelector('.rbac-badge')) return;

        const badge = document.createElement('span');
        badge.className = 'rbac-badge';
        badge.textContent = `Roles: ${roles.join(', ')}`;
        badge.style.cssText = `
            display: inline-block;
            background-color: #3b4151; /* Dark Grey */
            color: white;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 12px;
            font-weight: bold;
            margin-left: 10px;
            vertical-align: middle;
            line-height: normal;
        `;

        // Insert after path or description
        const pathData = block.querySelector('.opblock-summary-path');
        if (pathData) {
            pathData.appendChild(badge);
        } else if (summary) {
            summary.appendChild(badge);
        }
    };

    // Start
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initRBAC);
    } else {
        initRBAC();
    }
})();
