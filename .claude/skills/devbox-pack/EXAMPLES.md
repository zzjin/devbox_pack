# DevBox Pack Analyzer - Examples & Usage Patterns

This document provides comprehensive examples of how to use the DevBox Pack Analyzer AI skill for project analysis and execution plan generation.

## Example 1: React TypeScript Project

### Project Structure
```
my-react-app/
├── package.json
├── package-lock.json
├── tsconfig.json
├── src/
│   ├── App.tsx
│   ├── index.tsx
│   └── components/
├── public/
│   └── index.html
└── node_modules/
```

### Key Files Analysis
- **package.json**: Contains React dependencies and scripts
- **package-lock.json**: Indicates npm as package manager
- **tsconfig.json**: TypeScript configuration present
- **src/*.tsx**: TypeScript React files

### AI Analysis Process

#### Step 1: File System Scan
```
Detected Files:
✓ package.json (Node.js indicator: +40 points)
✓ package-lock.json (npm lock file: +20 points)
✓ tsconfig.json (TypeScript config: +5 points)
✓ node_modules/ (dependencies: +15 points)
✓ *.tsx files (TypeScript files: +10 points)
```

#### Step 2: Confidence Calculation
```
Total Score: 40 + 20 + 5 + 15 + 10 = 90/100
Confidence: 0.90 (High confidence)
```

#### Step 3: Framework Detection
```
package.json analysis:
- "react": "^18.2.0" → Framework: React
- "@types/react": "^18.2.0" → TypeScript support confirmed
```

#### Step 4: Package Manager Detection
```
Lock file analysis:
- package-lock.json exists → Package Manager: npm
```

### Generated Execution Plan
```json
{
  "provider": "node",
  "base": {
    "name": "node:20-slim",
    "platform": "linux/amd64"
  },
  "runtime": {
    "language": "javascript",
    "version": "20",
    "tools": ["npm"],
    "environment": {
      "NODE_ENV": "development"
    }
  },
  "apt": ["curl", "git"],
  "commands": {
    "dev": ["npm run dev"],
    "build": ["npm run build"],
    "start": ["npm start"]
  },
  "port": 3000,
  "evidence": {
    "files": ["package.json", "package-lock.json", "tsconfig.json"],
    "reason": "Detected Node.js project based on: package.json with React dependencies, npm lock file, TypeScript configuration"
  }
}
```

## Example 2: Django Python Project

### Project Structure
```
django-app/
├── requirements.txt
├── manage.py
├── myproject/
│   ├── __init__.py
│   ├── settings.py
│   ├── urls.py
│   └── wsgi.py
├── myapp/
│   ├── models.py
│   ├── views.py
│   └── urls.py
└── venv/
```

### AI Analysis Process

#### Step 1: File System Scan
```
Detected Files:
✓ requirements.txt (Python config: +40 points)
✓ manage.py (Django indicator: +20 points)
✓ *.py files (Python files: +20 points)
✓ venv/ (virtual environment: +10 points)
```

#### Step 2: Framework Detection
```
requirements.txt analysis:
- Django==4.2.0 → Framework: Django
- psycopg2-binary → Database: PostgreSQL
```

### Generated Execution Plan
```json
{
  "provider": "python",
  "base": {
    "name": "python:3.11-slim",
    "platform": "linux/amd64"
  },
  "runtime": {
    "language": "python",
    "version": "3.11",
    "tools": ["pip"],
    "environment": {
      "PYTHONPATH": "/app",
      "DJANGO_SETTINGS_MODULE": "myproject.settings"
    }
  },
  "apt": ["curl", "git", "build-essential", "libpq-dev"],
  "commands": {
    "install": ["pip install -r requirements.txt"],
    "dev": ["python manage.py runserver 0.0.0.0:8000"],
    "migrate": ["python manage.py migrate"],
    "start": ["python manage.py runserver 0.0.0.0:8000"]
  },
  "port": 8000,
  "evidence": {
    "files": ["requirements.txt", "manage.py"],
    "reason": "Detected Python Django project based on: requirements.txt with Django dependency, manage.py script"
  }
}
```

## Example 3: Spring Boot Java Project

### Project Structure
```
spring-app/
├── pom.xml
├── src/
│   └── main/
│       └── java/
│           └── com/
│               └── example/
│                   └── Application.java
└── target/
```

### AI Analysis Process

#### Step 1: File System Scan
```
Detected Files:
✓ pom.xml (Maven build file: +50 points)
✓ src/main/java/ (Java structure: +15 points)
✓ *.java files (Java files: +25 points)
```

#### Step 2: Framework Detection
```
pom.xml analysis:
- spring-boot-starter-web → Framework: Spring Boot
- spring-boot-starter-data-jpa → JPA support
```

### Generated Execution Plan
```json
{
  "provider": "java",
  "base": {
    "name": "openjdk:17-jdk-slim",
    "platform": "linux/amd64"
  },
  "runtime": {
    "language": "java",
    "version": "17",
    "tools": ["maven"],
    "environment": {
      "JAVA_OPTS": "-Xmx512m"
    }
  },
  "apt": ["curl", "git"],
  "commands": {
    "build": ["mvn clean package"],
    "dev": ["mvn spring-boot:run"],
    "start": ["java -jar target/*.jar"]
  },
  "port": 8080,
  "evidence": {
    "files": ["pom.xml", "src/main/java/"],
    "reason": "Detected Java Spring Boot project based on: Maven pom.xml with Spring Boot dependencies"
  }
}
```

## Example 4: Multi-Language Project (Node.js + Python)

### Project Structure
```
fullstack-app/
├── package.json          # Frontend (React)
├── requirements.txt      # Backend (FastAPI)
├── frontend/
│   ├── src/
│   └── public/
├── backend/
│   ├── main.py
│   └── api/
└── docker-compose.yml
```

### AI Analysis Process

#### Step 1: Multiple Provider Detection
```
Node.js Provider:
- package.json: +40 points
- Confidence: 0.40

Python Provider:
- requirements.txt: +40 points
- *.py files: +20 points
- Confidence: 0.60
```

#### Step 2: Best Match Selection
```
Best Match: Python (0.60 confidence)
Secondary: Node.js (0.40 confidence)
```

### Generated Execution Plan (Primary: Python)
```json
{
  "provider": "python",
  "base": {
    "name": "python:3.11-slim",
    "platform": "linux/amd64"
  },
  "runtime": {
    "language": "python",
    "version": "3.11",
    "tools": ["pip"],
    "environment": {
      "PYTHONPATH": "/app"
    }
  },
  "apt": ["curl", "git", "build-essential", "nodejs", "npm"],
  "commands": {
    "install": ["pip install -r requirements.txt", "npm install"],
    "dev": ["uvicorn main:app --host 0.0.0.0 --port 8000 --reload"],
    "build": ["npm run build"],
    "start": ["uvicorn main:app --host 0.0.0.0 --port 8000"]
  },
  "port": 8000,
  "evidence": {
    "files": ["requirements.txt", "main.py", "package.json"],
    "reason": "Detected Python FastAPI project with Node.js frontend based on: requirements.txt, FastAPI main.py, package.json"
  }
}
```

## Example 5: Static Website

### Project Structure
```
static-site/
├── index.html
├── styles.css
├── script.js
├── images/
│   ├── logo.png
│   └── hero.jpg
└── assets/
```

### AI Analysis Process

#### Step 1: File System Scan
```
Detected Files:
✓ index.html (HTML file: +20 points)
✓ *.css files (CSS files: +20 points)
✓ *.js files (JavaScript files: +10 points)
✓ images/ (static assets: +20 points)
```

#### Step 2: Confidence Calculation
```
Total Score: 20 + 20 + 10 + 20 = 70/100
Confidence: 0.70 (Medium-High confidence)
```

### Generated Execution Plan
```json
{
  "provider": "staticfile",
  "base": {
    "name": "nginx:alpine",
    "platform": "linux/amd64"
  },
  "runtime": {
    "language": "html",
    "version": "latest",
    "tools": ["nginx"],
    "environment": {}
  },
  "apt": [],
  "commands": {
    "start": ["nginx", "-g", "daemon off;"]
  },
  "port": 80,
  "evidence": {
    "files": ["index.html", "styles.css", "images/"],
    "reason": "Detected static website based on: HTML files, CSS stylesheets, static assets"
  }
}
```

## Example 6: Go Gin Web Service

### Project Structure
```
go-api/
├── go.mod
├── go.sum
├── main.go
├── handlers/
│   └── api.go
└── models/
    └── user.go
```

### AI Analysis Process

#### Step 1: File System Scan
```
Detected Files:
✓ go.mod (Go module: +50 points)
✓ go.sum (Go dependencies: +15 points)
✓ main.go (Go entry point: +5 points)
✓ *.go files (Go source: +30 points)
```

#### Step 2: Framework Detection
```
go.mod analysis:
- github.com/gin-gonic/gin v1.9.1 → Framework: Gin
```

### Generated Execution Plan
```json
{
  "provider": "go",
  "base": {
    "name": "golang:1.21-alpine",
    "platform": "linux/amd64"
  },
  "runtime": {
    "language": "go",
    "version": "1.21",
    "tools": ["go"],
    "environment": {
      "CGO_ENABLED": "0",
      "GOOS": "linux"
    }
  },
  "apt": ["curl", "git", "ca-certificates"],
  "commands": {
    "build": ["go build -o app ."],
    "dev": ["go run main.go"],
    "start": ["./app"]
  },
  "port": 8080,
  "evidence": {
    "files": ["go.mod", "main.go"],
    "reason": "Detected Go Gin web service based on: go.mod with Gin framework dependency"
  }
}
```

## Usage Patterns for AI Assistants

### Pattern 1: Basic Project Analysis
```markdown
1. List project files using file system tools
2. Apply detection logic for each provider
3. Calculate confidence scores
4. Select best match
5. Generate execution plan
```

### Pattern 2: Framework-Specific Analysis
```markdown
1. Identify primary language (highest confidence)
2. Examine configuration files for framework signatures
3. Apply framework-specific detection rules
4. Generate framework-optimized execution plan
```

### Pattern 3: Multi-Language Project Handling
```markdown
1. Detect all languages with confidence > 0.3
2. Rank by confidence scores
3. Select primary language as base
4. Include secondary language tools in apt packages
5. Generate hybrid execution plan
```

### Pattern 4: Version Detection Strategy
```markdown
1. Check version-specific files (.nvmrc, .python-version)
2. Parse configuration files (package.json engines, pyproject.toml)
3. Examine lock files for version constraints
4. Apply sensible defaults if no version found
```

### Pattern 5: Evidence Collection
```markdown
1. Document all matching indicator files
2. Explain confidence score breakdown
3. Provide reasoning for framework selection
4. Include metadata about project characteristics
```

## Common Analysis Scenarios

### Scenario 1: Ambiguous Projects
When multiple providers have similar confidence scores:
- Use provider priority as tiebreaker
- Consider file count and project structure
- Prefer more specific over generic providers

### Scenario 2: Legacy Projects
For older projects without modern configuration:
- Rely on file extensions and directory structure
- Use lower confidence thresholds
- Apply conservative base image selection

### Scenario 3: Monorepo Projects
For projects with multiple sub-projects:
- Analyze each subdirectory separately
- Aggregate results with weighted confidence
- Generate multi-service execution plans

### Scenario 4: Configuration-Heavy Projects
For projects with extensive configuration:
- Parse configuration files for framework hints
- Consider build tool configurations
- Include configuration-specific packages

This comprehensive example set demonstrates how the AI skill analyzes various project types and generates appropriate execution plans through structured reasoning and decision trees.