import fs from 'fs'
import path from 'path'

const CurrentDir = process.cwd();

let SCAN_DIR = path.join(CurrentDir, 'src/routes');
if (!fs.existsSync(SCAN_DIR)) {
    SCAN_DIR = path.join(CurrentDir, 'web/src/routes');
}

console.log(`[i18n-scanner] Scanning: ${SCAN_DIR}`);

const IGNORE_PATTERNS = [
    /^\s*$/,
    /^\d+$/,
    /^[!@#$%^&*()_+=\-{}[\]|\\:;"'<>,.?/`~]+$/, 
    /^{.*}$/,
    /^&[a-z]+;$/,
];

const IGNORE_FILES = [
    '.gen.ts',
    '.d.ts'
];

// Heuristic regex
const JSX_TEXT_REGEX = />([^<{]+)</g
const ATTRIBUTE_REGEX = /\b(placeholder|title|alt|aria-label)="([^"]+)"/g

function getAllFiles(dirPath: string, arrayOfFiles: string[] = []) {
    if (!fs.existsSync(dirPath)) {
        console.error(`Directory not found: ${dirPath}`);
        return [];
    }
    const files = fs.readdirSync(dirPath)

    files.forEach((file) => {
        const fullPath = path.join(dirPath, file);
        if (fs.statSync(fullPath).isDirectory()) {
            getAllFiles(fullPath, arrayOfFiles)
        } else {
            if ((file.endsWith('.tsx') || file.endsWith('.ts')) && !IGNORE_FILES.some(ignore => file.endsWith(ignore))) {
                arrayOfFiles.push(fullPath)
            }
        }
    })

    return arrayOfFiles
}

function run() {
    const files = getAllFiles(SCAN_DIR);
    console.log(`[i18n-scanner] Found ${files.length} files.`);

    let totalErrors = 0;

    files.forEach(file => {
        const content = fs.readFileSync(file, 'utf-8');
        const errors: string[] = [];

        let match;


        while ((match = JSX_TEXT_REGEX.exec(content)) !== null) {
            const text = match[1].trim();
            if (text && !IGNORE_PATTERNS.some(p => p.test(text))) {
                errors.push(`Line: ? | Text: "${text}"`);
            }
        }

        while ((match = ATTRIBUTE_REGEX.exec(content)) !== null) {
            const attr = match[1];
            const value = match[2].trim();
            if (value && !IGNORE_PATTERNS.some(p => p.test(value))) {
                errors.push(`Line: ? | Attribute [${attr}]: "${value}"`);
            }
        }

        if (errors.length > 0) {
            console.log(`\n❌ Potential unlocalized strings in: ${path.relative(CurrentDir, file)}`);
            errors.forEach(e => console.log(`   ${e}`));
            totalErrors += errors.length;
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
