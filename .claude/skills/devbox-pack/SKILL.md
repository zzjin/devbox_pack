---
name: analyzing-projects
description: Analyze project structure to detect languages, frameworks, and generate containerized execution plans. Supports Node.js, Python, Java, Go, PHP, Ruby, Deno, Rust, Shell, and static files.
---

# DevBox Pack Analyzer

Analyze project structure to detect languages, frameworks, and generate containerized execution plans.

> **📖 For detailed examples and usage patterns, see [EXAMPLES.md](EXAMPLES.md)**

## Core Functionality

1. **Project Detection** - Scan files to identify language, framework, and tools
2. **Confidence Scoring** - Calculate match confidence based on file indicators  
3. **Execution Plan Generation** - Create containerized development configuration

## Supported Languages

**Node.js**, **Python**, **Java**, **Go**, **PHP**, **Ruby**, **Deno**, **Rust**, **Shell**, **Static Files**

## JSON Output Structure

Generate output matching the exact Go implementation structure:

### ExecutionPlan
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
    "dev": ["npm install", "npm run dev -- --host 0.0.0.0 --port ${PORT}"],
    "build": ["npm install", "npm run build"],
    "start": ["npm start -- --host 0.0.0.0 --port ${PORT}"]
  },
  "port": 3000,
  "evidence": {
    "files": ["package.json", "package-lock.json"],
    "reason": "Detected Node.js project based on package.json and npm lock file"
  }
}
```

### DetectResult
```json
{
  "matched": true,
  "provider": "node",
  "confidence": 0.90,
  "evidence": {
    "files": ["package.json", "package-lock.json", "tsconfig.json"],
    "reason": "Node.js project with TypeScript configuration"
  },
  "language": "javascript",
  "framework": "react",
  "version": "18.2.0",
  "packageManager": {
    "name": "npm",
    "lockFile": "package-lock.json",
    "configFile": "package.json",
    "useCorepack": false
  },
  "buildTools": ["webpack", "typescript"],
  "metadata": {
    "hasTypeScript": true,
    "scripts": ["dev", "build", "start"]
  }
}
```

## Key Structure Requirements

- **provider**: Language identifier (node, python, java, go, php, ruby, deno, rust, shell, staticfile)
- **base**: Object with `name` and optional `platform`
- **runtime**: Object with `language`, optional `version`, `tools` array, optional `environment` map
- **apt**: Array of system packages
- **commands**: Object with optional `dev`, `build`, `start` arrays
- **port**: Integer port number
- **evidence**: Object with `files` array and `reason` string
- **packageManager**: Object with `name`, optional `lockFile`, `configFile`, `useCorepack`

## Detection Workflow (4 Phases)

**Phase 1: Source Code Acquisition**
- Scan project files and directory structure
- Collect file indicators and patterns

**Phase 2: Provider Detection** 
- Execute detection logic for each provider by priority order
- Calculate confidence scores using weighted indicators
- Select provider with highest confidence (threshold: 0.3)

**Phase 3: Plan Generation**
- Analyze dependencies and package managers
- Generate development, build, and production commands
- Determine base image and system requirements

**Phase 4: Output Generation**
- Format execution plan with detection evidence
- Include confidence metrics and metadata

## Provider Priority Order

Lower number = higher priority:
1. **Static File** (10) - HTML/CSS/JS static sites
2. **Shell** (30) - Shell scripts and executables  
3. **PHP** (60) - PHP applications
4. **Ruby** (65) - Ruby/Rails applications
5. **Java** (70) - Java/Spring applications
6. **Go** (75) - Go applications
7. **Node.js** (80) - JavaScript/TypeScript applications
8. **Python** (80) - Python applications
9. **Rust** (85) - Rust applications
10. **Deno** (95) - Deno applications

## Critical Networking Requirements

**MANDATORY 0.0.0.0 Binding:**
- All commands MUST bind to `0.0.0.0` (never localhost/127.0.0.1)
- Use `${PORT}` environment variable with framework defaults
- Required for container compatibility

**Command Structure:**
- **dev**: [install_dependencies, start_dev_server_with_0.0.0.0_binding]
- **build**: [install_dependencies, run_build_command]  
- **start**: [production_start_command_with_0.0.0.0_binding]

**Language-Specific 0.0.0.0 Binding Examples:**
- **Node.js**: `npm run dev -- --host 0.0.0.0 --port ${PORT}`
- **Python Django**: `python manage.py runserver 0.0.0.0:8000`
- **Python FastAPI**: `uvicorn main:app --host 0.0.0.0 --port ${PORT}`
- **PHP**: `php -S 0.0.0.0:8000`
- **Ruby Rails**: `bundle exec rails server -b 0.0.0.0 -p ${PORT}`
- **Ruby Sinatra**: `bundle exec rackup -o 0.0.0.0 -p ${PORT}`