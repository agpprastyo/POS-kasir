
(function () {
    const POLL_INTERVAL = 500;
    const MAX_RETRIES = 20;
    let attempts = 0;

    const observer = new MutationObserver((mutations) => {
        enrichUI();
        hideRawExtensions();
    });

    const initRBAC = () => {
        // Wait until UI system is ready
        if (typeof window.ui === 'undefined' || typeof window.ui.spec !== 'function') {
            if (attempts < MAX_RETRIES) {
                attempts++;
                setTimeout(initRBAC, POLL_INTERVAL);
            }
            return;
        }

        const app = document.getElementById('swagger-ui');
        if (app) {
            observer.observe(app, {
                childList: true,
                subtree: true
            });
            // Initial attempt
            setTimeout(() => {
                enrichUI();
                hideRawExtensions();
            }, 1000);
        }
    };

    const hideRawExtensions = () => {
        // Find all table cells that might contain the extension key
        const cells = document.querySelectorAll('td');
        cells.forEach(td => {
            if (td.textContent.trim() === 'x-roles') {
                const row = td.closest('tr');
                // Hide the row if found
                if (row) {
                    row.style.display = 'none';

                    // Also check if the table becomes empty, if so hide the table/headers
                    const table = row.closest('table');
                    if (table) {
                        const visibleRows = Array.from(table.querySelectorAll('tr')).filter(r => r.style.display !== 'none');
                        if (visibleRows.length === 0) {
                            table.style.display = 'none';
                            // Try to hide the "Extensions" header if it exists nearby
                            const wrapper = table.closest('.opblock-section-header');
                            if (wrapper) wrapper.style.display = 'none';
                        }
                    }
                }
            }
        });
    };

    const enrichUI = () => {
        try {
            // Get spec from Swagger UI system
            let spec = null;
            if (window.ui.spec().selectors && window.ui.spec().selectors.specJson) {
                spec = window.ui.spec().selectors.specJson().toJS();
            } else if (window.ui.getSystem().specSelectors && window.ui.getSystem().specSelectors.specJson) {
                spec = window.ui.getSystem().specSelectors.specJson().toJS();
            }

            if (!spec || !spec.paths) return;

            const opBlocks = document.querySelectorAll('.opblock');

            opBlocks.forEach(block => {
                // Don't skip if already enriched, because re-rendering might remove the badge 
                // but keep the dataset. Instead check if badge exists.
                if (block.querySelector('.rbac-badge')) return;

                const methodEl = block.querySelector('.opblock-summary-method');
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
            });
        } catch (e) {
            console.warn('RBAC Enrichment error:', e);
        }
    };

    const createBadge = (block, roles) => {
        const summary = block.querySelector('.opblock-summary-operation-id') || block.querySelector('.opblock-summary-description') || block.querySelector('.opblock-summary');

        if (!summary) return;

        // Prevent duplicate badges
        if (block.querySelector('.rbac-badge')) return;

        const badge = document.createElement('span');
        badge.className = 'rbac-badge';

        // Format roles nicely
        const rolesText = roles.map(r => r.toUpperCase()).join(', ');
        badge.textContent = rolesText; // Just the roles, e.g. "ADMIN, CASHIER"

        // Create a wrapper for better positioning
        const wrapper = document.createElement('div');
        wrapper.className = 'rbac-badge-wrapper';
        wrapper.style.cssText = `
            display: inline-flex;
            align-items: center;
            margin-left: 10px;
            vertical-align: middle;
        `;

        // Style the badge itself
        badge.style.cssText = `
            display: inline-block;
            background-color: #49cc90; /* Green (GET) style or custom */
            color: white;
            padding: 4px 10px;
            border-radius: 4px;
            font-size: 11px;
            font-weight: bold;
            text-transform: uppercase;
            box-shadow: 0 1px 2px rgba(0,0,0,0.1);
        `;

        // Different colors for different roles?
        if (roles.includes('admin')) {
            badge.style.backgroundColor = '#f93e3e'; // Red for admin
        } else if (roles.includes('manager')) {
            badge.style.backgroundColor = '#fca130'; // Orange for manager
        } else {
            badge.style.backgroundColor = '#49cc90'; // Green for others
        }

        // Add label prefix
        const label = document.createElement('span');
        label.textContent = 'ROLES: ';
        label.style.cssText = `
            font-size: 10px;
            font-weight: bold;
            color: #888;
            margin-right: 4px;
        `;

        wrapper.appendChild(label);
        wrapper.appendChild(badge);

        // Insert after path
        const pathData = block.querySelector('.opblock-summary-path');
        if (pathData) {
            // Insert after the path link
            pathData.parentNode.insertBefore(wrapper, pathData.nextSibling);
        } else {
            summary.appendChild(wrapper);
        }
    };

    // Start
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initRBAC);
    } else {
        initRBAC();
    }
})();
