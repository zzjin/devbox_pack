# Static Files Examples

Complete examples of static website analysis and deployment configuration.

## Example 1: Pure HTML/CSS/JS Website

**User request**: "Pack this HTML site"

**Detected files**:
- `index.html`
- `css/style.css`
- `js/app.js`
- No backend frameworks
- No package.json

**Output**:
```json
[
  {
    "language": "static",
    "apt": [],
    "dev": {
      "environment": {},
      "commands": ["python3 -m http.server 8080 --bind 0.0.0.0"]
    },
    "prod": {
      "environment": {},
      "commands": ["python3 -m http.server 8080 --bind 0.0.0.0"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\npython3 -m http.server 8080 --bind 0.0.0.0",
    "evidence": {
      "files": ["index.html", "css/style.css", "js/app.js"],
      "reason": "Static HTML/CSS/JS website with no backend framework"
    }
  }
]
```

**Key Points**:
- Lowest priority (#9) in language detection
- Only detected when NO backend framework present
- Same command for dev and prod (simple HTTP server)
- No build step needed
- No setup commands
- No version field (not applicable)

---

## Example 2: Single Page Application (Built)

**User request**: "Serve built SPA"

**Detected files**:
- `dist/` or `build/` directory with `index.html`
- `dist/assets/` with bundled JS/CSS
- No source files (already built)
- No package.json in root

**Output**:
```json
[
  {
    "language": "static",
    "apt": [],
    "dev": {
      "environment": {},
      "commands": ["python3 -m http.server 8080 --bind 0.0.0.0 --directory dist"]
    },
    "prod": {
      "environment": {},
      "commands": ["python3 -m http.server 8080 --bind 0.0.0.0 --directory dist"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\npython3 -m http.server 8080 --bind 0.0.0.0 --directory dist",
    "evidence": {
      "files": ["dist/index.html", "dist/assets/"],
      "reason": "Pre-built static SPA detected in dist directory"
    }
  }
]
```

**Key Points**:
- Serve from `dist/` or `build/` directory
- Assumes build already completed (CI/CD pipeline)
- Common for production-ready SPAs

---

## Example 3: Static Site with Multiple Pages

**User request**: "Deploy static documentation site"

**Detected files**:
- `index.html`
- `about.html`, `contact.html`
- `assets/` directory
- No dynamic backend

**Output**:
```json
[
  {
    "language": "static",
    "apt": [],
    "dev": {
      "environment": {},
      "commands": ["python3 -m http.server 8080 --bind 0.0.0.0"]
    },
    "prod": {
      "environment": {},
      "commands": ["python3 -m http.server 8080 --bind 0.0.0.0"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\npython3 -m http.server 8080 --bind 0.0.0.0",
    "evidence": {
      "files": ["index.html", "about.html", "contact.html", "assets/"],
      "reason": "Multi-page static website detected"
    }
  }
]
```

---

## Static Files Detection Patterns

### Detection Criteria

**MUST have**:
- `index.html` file in root or subdirectory

**MUST NOT have** (any backend framework):
- `package.json` (Node.js)
- `requirements.txt` (Python)
- `go.mod` (Go)
- `Gemfile` (Ruby)
- `composer.json` (PHP)
- `Cargo.toml` (Rust)
- `pom.xml` or `build.gradle` (Java)
- `deno.json` (Deno)

**Priority**: #9 (second-lowest, only above Shell scripts)

### Common File Structures

**Traditional website**:
```
website/
├── index.html
├── about.html
├── css/
│   └── style.css
├── js/
│   └── script.js
└── images/
    └── logo.png
```

**Built SPA**:
```
project/
└── dist/
    ├── index.html
    ├── assets/
    │   ├── index-[hash].js
    │   └── index-[hash].css
    └── favicon.ico
```

**Documentation site**:
```
docs/
├── index.html
├── getting-started.html
├── api/
│   └── reference.html
└── assets/
```

### Static File Server Options

**Python 3** (recommended - always available):
```bash
python3 -m http.server 8080 --bind 0.0.0.0
python3 -m http.server 8080 --bind 0.0.0.0 --directory dist
```

**Alternative servers** (if available):
```bash
# Node.js serve package
npx serve -l 8080 dist

# PHP built-in server
php -S 0.0.0.0:8080

# Caddy
caddy file-server --listen 0.0.0.0:8080
```

**Production servers** (for reference, not used in config):
- nginx
- Apache httpd
- Caddy
- Lighttpd

### Default Port
**8080** (standard for static file servers)

### No Build Step
Static sites detected by this skill are assumed to be:
- Already built (if SPA)
- Or require no build (plain HTML)

If build is needed, project should have `package.json` and be detected as Node.js.

### No Environment Variables
Static sites typically need no environment configuration.

Exception: SPAs may use build-time environment variables, but these are baked into the built files.

### No APT Dependencies
Python 3 is available in all base images, so `[]` for apt.

### Production Considerations

**For Production Deployment**:
1. **Use a real web server**: nginx, Caddy, or CDN
2. **Enable GZIP/Brotli**: Compression for faster loading
3. **Set cache headers**: Leverage browser caching
4. **HTTPS**: Always use TLS in production
5. **CDN**: Consider CloudFront, Cloudflare, Vercel for global distribution

**Python http.server limitations**:
- Single-threaded
- No compression
- No caching headers
- Not optimized for production

**However**: For simple static sites or internal tools, Python's http.server is sufficient.

### SPA Routing

**Problem**: SPAs with client-side routing need all routes to serve `index.html`

**Solution** (not in basic config):
```bash
# For nginx (production):
location / {
    try_files $uri $uri/ /index.html;
}

# For Python http.server: Limited support
# Better to use nginx or custom server
```

### Content Types
Python's http.server automatically sets correct MIME types:
- `.html` → `text/html`
- `.css` → `text/css`
- `.js` → `application/javascript`
- `.json` → `application/json`
- `.png` → `image/png`
- etc.

### Directory Listing
By default, `http.server` shows directory listing if no `index.html`.

Disable in production with proper web server config.

### Common Use Cases
1. **Landing pages**: Simple marketing sites
2. **Documentation**: Generated docs (Sphinx, MkDocs output)
3. **Built SPAs**: React/Vue/Angular production builds
4. **Prototypes**: Quick HTML mockups
5. **Static blogs**: Jekyll, Hugo output

### Not Considered Static
These have backend frameworks and are detected differently:
- **Next.js** → Detected as Node.js (has package.json)
- **Gatsby** → Detected as Node.js (has package.json)
- **Jekyll** → Would need Ruby (has Gemfile)
- **Hugo** → Would need Go (has go.mod) or detected as static if only built output
- **Eleventy** → Detected as Node.js (has package.json)

### When to Use Static Detection
- Pre-built production assets only
- No source code to compile
- No package manager configuration
- Pure HTML/CSS/JS files

### Security Considerations
1. **Path traversal**: Python http.server prevents `../` attacks
2. **No server-side code**: Can't execute malicious scripts
3. **Read-only**: Server only serves files, no writes
4. **CORS**: Python http.server has no CORS headers by default

### Performance
- **Fast startup**: No compilation or setup
- **Low memory**: Minimal resource usage
- **Scalability**: For high traffic, use nginx/CDN
- **Caching**: No automatic cache control

### Example: Documentation Site

**Scenario**: MkDocs-generated site in `site/` directory

```json
{
  "language": "static",
  "apt": [],
  "dev": {
    "commands": ["python3 -m http.server 8080 --bind 0.0.0.0 --directory site"]
  },
  "prod": {
    "commands": ["python3 -m http.server 8080 --bind 0.0.0.0 --directory site"]
  },
  "port": 8080,
  "entrypoint": "#!/bin/bash\n\npython3 -m http.server 8080 --bind 0.0.0.0 --directory site"
}
```
