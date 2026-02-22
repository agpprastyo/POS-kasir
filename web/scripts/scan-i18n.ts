import fs from 'fs'
import path from 'path'

const CurrentDir = process.cwd();

let SCAN_DIR = path.join(CurrentDir, 'src');
if (!fs.existsSync(SCAN_DIR)) {
    SCAN_DIR = path.join(CurrentDir, 'web/src');
}

let LOCALES_DIR = path.join(CurrentDir, 'src/lib/locales');
if (!fs.existsSync(LOCALES_DIR)) {
    LOCALES_DIR = path.join(CurrentDir, 'web/src/lib/locales');
}

console.log(`[i18n-scanner] Scanning: ${SCAN_DIR}`);

const IGNORE_PATTERNS = [
    /^\s*$/,
    /^\d+$/,
    /^[!@#$%^&*()_+=\-{}[\]|\\:;"'<>,.?/`~]+$/,
    /^{.*}$/,
    /^&[a-z]+;$/,
    // Common code patterns
    /^[A-Z_]+$/, // ENUMS
    /^[a-z]+[A-Z][a-zA-Z]+$/, // camelCase (like setSomething, handleAction, etc)
    /^\(null\)/, // React useState initializations
    /===|!==|&&|\|\||>=|<=|=>|\btypeof\b|\binstanceof\b|function\s*\(|console\.log/,
    /^import\s+/,
    /^export\s+/,
    /className="/,
    /\]\s*=\s*useState/,
    /\]\s*=\s*use/,
    /}\s*function/i,
    /^\)\s*\}/,
    /function\s+[A-Za-z0-9_]+\s*\(/,
    /\bnavigate\(/,
    /\bdispatch\(/,
    /\.mutate\(/,
    /\.handleChange\(/,
    /\.parse\(/,
    /\.ensureQueryData\(/,
    /^as\s+[A-Z]/,
    /\]\}\s*children=/,
    /^\)\s*:\s*\(/,
    /t\('[^']+'\)/,
    /^void$/,
    /^\)\s*\}\s*return\b/,
    /\bformatDate\(/,
    /\bformatCurrency\(/,
    /labelFormatter=/,
    /search\.page/,
    /search\.start_date/,
    /search\.end_date/,
    /^\s*:\s*null\s*$/,
    /^\s*:\s*undefined\s*$/,
    /^\s*:\s*(true|false)\s*$/,
    /=\s*$/,
    /=>/,
    /\b(ReactNode|null|undefined)\b/,
    /^\s*const\s+/,
    /^\s*let\s+/,
    /^\s*var\s+/,
    /^\s*import\s+/,
    /^\s*export\s+/,
    /^\s*interface\s+/,
    /^\s*type\s+/,
    /^\s*return\s+/,
    /^\s*return\s*\(/,
    /}\s*else\b/,
    /^\s*if\s*\(/,
    /^\s*;/,
    /^[A-Za-z0-9_]+\s*\|\s*null\b/,
    /use[A-Z][A-Za-z0-9_]+\b/,
    /^\s*\)\s*;\s*$/,
    /^\s*\}\s*;\s*$/,
    /\};?\s+return\s+\(/,
    /\[\s*state\.canSubmit/,
    /\buseRef\b/,
    /\b[a-zA-Z0-9_]+\s*:\s*[A-Z][a-zA-Z0-9_]+\s*\|\s*null/,
    /\bhandle[A-Z][a-zA-Z0-9_]+\s*\(/,
    /\btrigger[A-Z][a-zA-Z0-9_]+\s*\(/,
    /\bon[A-Z][a-zA-Z0-9_]+\s*\(/,
    /search,\s*loader:\s*\(/,
    /\]\}\s*children=/,
    /^`Rp\$$/,
    /^`Rp\$$/,
    /\bset[A-Z][a-zA-Z0-9_]+\s*\(/,
    /\bopen[A-Z][a-zA-Z0-9_]+\s*\(/,
    /\badd[A-Z][a-zA-Z0-9_]+\s*\(/,
    /\bfield\.[a-zA-Z0-9_]+\s*\(/,
    /\?\s*['"][a-zA-Z0-9_\-]+['"]\s*:\s*['"][a-zA-Z0-9_\-]+['"]/,
    /^\s*0\s*;\s*$/,
    /^\s*return\s*\(\s*$/,

];

const IGNORE_FILES = [
    '.gen.ts',
    '.d.ts',
];

// Heuristic regex
const JSX_TEXT_REGEX = />([^<{]+)[<\{]/g
const ATTRIBUTE_REGEX = /\b(placeholder|title|alt|aria-label)="([^"]+)"/g
const TOAST_REGEX = /\btoast(?:\.(?:success|error|warning|info))?\(\s*(['"`])(.*?)\1/g
const TRANSLATION_KEY_REGEX = /\bt\(\s*['"]([^'"]+)['"]/g

function getAllFiles(dirPath: string, arrayOfFiles: string[] = []) {
    if (!fs.existsSync(dirPath)) {
        console.error(`Directory not found: ${dirPath}`);
        return [];
    }
    const files = fs.readdirSync(dirPath)

    files.forEach((file) => {
        const fullPath = path.join(dirPath, file);
        if (fs.statSync(fullPath).isDirectory()) {
            // Also ignore Shadcn UI components folder
            if (file !== 'locales' && file !== 'generated' && !fullPath.includes('components/ui')) {
                getAllFiles(fullPath, arrayOfFiles)
            }
        } else {
            if ((file.endsWith('.tsx') || file.endsWith('.ts')) && !IGNORE_FILES.some(ignore => file.endsWith(ignore))) {
                arrayOfFiles.push(fullPath)
            }
        }
    })

    return arrayOfFiles
}


function flattenObject(ob: any): Record<string, string> {
    var toReturn: Record<string, string> = {};
    for (var i in ob) {
        if (!ob.hasOwnProperty(i)) continue;
        if ((typeof ob[i]) == 'object' && ob[i] !== null) {
            var flatObject = flattenObject(ob[i]);
            for (var x in flatObject) {
                if (!flatObject.hasOwnProperty(x)) continue;
                toReturn[i + '.' + x] = flatObject[x];
            }
        } else {
            toReturn[i] = ob[i];
        }
    }
    return toReturn;
}

function loadLocales(): { id: Record<string, string>, en: Record<string, string> } {
    let id: Record<string, string> = {};
    let en: Record<string, string> = {};
    try {
        const idPath = path.join(LOCALES_DIR, 'id.json');
        const enPath = path.join(LOCALES_DIR, 'en.json');
        if (fs.existsSync(idPath)) id = flattenObject(JSON.parse(fs.readFileSync(idPath, 'utf8')));
        if (fs.existsSync(enPath)) en = flattenObject(JSON.parse(fs.readFileSync(enPath, 'utf8')));
    } catch (e) {
        console.warn("Failed to load locales", e);
    }
    return { id, en };
}

function run() {
    const files = getAllFiles(SCAN_DIR);
    console.log(`[i18n-scanner] Found ${files.length} files.`);

    const locales = loadLocales();

    let totalErrors = 0;

    files.forEach(file => {
        const content = fs.readFileSync(file, 'utf-8');
        const errors: string[] = [];

        let match;

        if (file.endsWith('.tsx')) {
            while ((match = JSX_TEXT_REGEX.exec(content)) !== null) {
                let text = match[1].trim();
                // Strip trailing parentheses or weird symbols captured
                text = text.replace(/[\)\]}]+$/, '').trim();

                // must contain a letter and not match code patterns
                if (text && /[a-zA-Z]/.test(text) && !IGNORE_PATTERNS.some(p => p.test(text))) {
                    errors.push(`Hardcoded UI Text: "${text}"`);
                }
            }

            while ((match = ATTRIBUTE_REGEX.exec(content)) !== null) {
                const attr = match[1];
                const value = match[2].trim();
                if (value && /[a-zA-Z]/.test(value) && !IGNORE_PATTERNS.some(p => p.test(value))) {
                    errors.push(`Hardcoded Attribute [${attr}]: "${value}"`);
                }
            }
        }

        // Toasts can be in both .ts and .tsx files
        while ((match = TOAST_REGEX.exec(content)) !== null) {
            const value = match[2].trim();
            if (value && /[a-zA-Z]/.test(value) && !IGNORE_PATTERNS.some(p => p.test(value))) {
                errors.push(`Hardcoded Toast: "${value}"`);
            }
        }

        // Check for missing translations
        while ((match = TRANSLATION_KEY_REGEX.exec(content)) !== null) {
            const key = match[1];
            if (!locales.id[key]) {
                errors.push(`Missing translation key '${key}' in id.json`);
            }
            if (!locales.en[key]) {
                errors.push(`Missing translation key '${key}' in en.json`);
            }
        }

        if (errors.length > 0) {
            console.log(`\n❌ Issues in: ${path.relative(CurrentDir, file)}`);
            // Deduplicate errors for the same file
            const uniqueErrors = Array.from(new Set(errors));
            uniqueErrors.forEach(e => console.log(`   ${e}`));
            totalErrors += uniqueErrors.length;
        }
    });

    console.log(`\n---------------------------------------------------`);
    if (totalErrors > 0) {
        console.log(`[i18n-scanner] Finished. Found ${totalErrors} potential issues.`);
        process.exit(1);
    } else {
        console.log(`[i18n-scanner] Finished. No obvious issues found.`);
        process.exit(0);
    }
}

run();
