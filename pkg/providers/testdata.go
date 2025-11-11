package providers

// TestData provides common test data templates for provider tests
type TestData struct {
	Files map[string]string
}

// Common test data templates
var (
	// GoTestData provides Go project test data
	GoTestData = TestData{
		Files: map[string]string{
			"go.mod": `module example.com/myapp

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/stretchr/testify v1.8.4
)`,
			"main.go": `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
		},
	}

	// NodeTestData provides Node.js project test data
	NodeTestData = TestData{
		Files: map[string]string{
			"package.json": `{
	"name": "test-project",
	"version": "1.0.0",
	"scripts": {
		"start": "node index.js",
		"build": "webpack"
	},
	"dependencies": {
		"react": "^18.0.0",
		"react-dom": "^18.0.0"
	},
	"engines": {
		"node": ">=18.0.0"
	}
}`,
			"index.js": `console.log("Hello, World!");`,
		},
	}

	// PythonTestData provides Python project test data
	PythonTestData = TestData{
		Files: map[string]string{
			"requirements.txt": `Django==4.2.0
requests==2.31.0
pytest==7.4.0
gunicorn==21.2.0`,
			"app.py": `from flask import Flask
app = Flask(__name__)

@app.route('/')
def hello():
    return "Hello, World!"

if __name__ == '__main__':
    app.run()`,
		},
	}

	// JavaMavenTestData provides Java Maven project test data
	JavaMavenTestData = TestData{
		Files: map[string]string{
			"pom.xml": `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>1.0.0</version>
    <properties>
        <maven.compiler.source>17</maven.compiler.source>
        <maven.compiler.target>17</maven.compiler.target>
    </properties>
</project>`,
			"src/main/java/Main.java": `public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}`,
		},
	}

	// JavaGradleTestData provides Java Gradle project test data
	JavaGradleTestData = TestData{
		Files: map[string]string{
			"build.gradle": `plugins {
    id 'java'
}

java {
    sourceCompatibility = JavaVersion.VERSION_17
    targetCompatibility = JavaVersion.VERSION_17
}`,
			"src/main/java/Main.java": `public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}`,
		},
	}

	// RustTestData provides Rust project test data
	RustTestData = TestData{
		Files: map[string]string{
			"Cargo.toml": `[package]
name = "hello-world"
version = "0.1.0"
edition = "2021"

[dependencies]
serde = { version = "1.0", features = ["derive"] }
tokio = { version = "1.0", features = ["full"] }
`,
			"src/main.rs": `fn main() {
    println!("Hello, World!");
}`,
		},
	}

	// PHPTestData provides PHP project test data
	PHPTestData = TestData{
		Files: map[string]string{
			"composer.json": `{
    "name": "example/my-app",
    "description": "A sample PHP application",
    "type": "project",
    "require": {
        "php": "^8.1",
        "laravel/framework": "^10.0"
    },
    "require-dev": {
        "phpunit/phpunit": "^10.0"
    },
    "autoload": {
        "psr-4": {
            "App\\": "app/"
        }
    }
}`,
			"index.php": `<?php
echo "Hello, World!";
?>`,
		},
	}

	// RubyTestData provides Ruby project test data
	RubyTestData = TestData{
		Files: map[string]string{
			"Gemfile": `source 'https://rubygems.org'

ruby '3.1.0'

gem 'rails', '~> 7.0.0'
gem 'sqlite3', '~> 1.4'
gem 'puma', '~> 5.0'
gem 'bootsnap', '>= 1.4.4', require: false

group :development, :test do
  gem 'byebug', platforms: [:mri, :mingw, :x64_mingw]
  gem 'rspec-rails'
end
`,
			"app.rb": `puts "Hello, World!"`,
		},
	}

	// DenoTestData provides Deno project test data
	DenoTestData = TestData{
		Files: map[string]string{
			"deno.json": `{
		"tasks": {
			"dev": "deno run --watch main.ts",
			"start": "deno run main.ts"
		},
		"imports": {
			"std/": "https://deno.land/std@0.200.0/"
		}
	}`,
			"main.ts": `console.log("Hello, World!");`,
		},
	}

	// ShellTestData provides Shell project test data
	ShellTestData = TestData{
		Files: map[string]string{
			"install.sh": `#!/bin/bash
echo "Installing..."
echo "Installation complete!"
`,
			"setup.sh": `#!/bin/bash
echo "Setting up..."
echo "Setup complete!"
`,
		},
	}

	// StaticFileTestData provides static file project test data
	StaticFileTestData = TestData{
		Files: map[string]string{
			"index.html": `<!DOCTYPE html>
<html>
<head>
    <title>Hello World</title>
</head>
<body>
    <h1>Hello, World!</h1>
</body>
</html>`,
			"styles.css": `body {
    font-family: Arial, sans-serif;
    margin: 0;
    padding: 20px;
}`,
			"script.js": `console.log("Hello, World!");`,
		},
	}
)

// Framework-specific test data
var (
	// NodeReactTestData provides Node.js React project test data
	NodeReactTestData = TestData{
		Files: map[string]string{
			"package.json": `{
	"name": "react-app",
	"version": "1.0.0",
	"dependencies": {
		"react": "^18.0.0",
		"react-dom": "^18.0.0"
	},
	"scripts": {
		"start": "react-scripts start",
		"build": "react-scripts build"
	}
}`,
			"src/App.js": `import React from 'react';

function App() {
  return (
    <div className="App">
      <h1>Hello, React!</h1>
    </div>
  );
}

export default App;`,
		},
	}

	// PythonDjangoTestData provides Python Django project test data
	PythonDjangoTestData = TestData{
		Files: map[string]string{
			"requirements.txt": `Django==4.2.0
psycopg2-binary==2.9.7
gunicorn==21.2.0`,
			"manage.py": `#!/usr/bin/env python
"""Django's command-line utility for administrative tasks."""
import os
import sys

if __name__ == '__main__':
    os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'myproject.settings')
    try:
        from django.core.management import execute_from_command_line
    except ImportError as exc:
        raise ImportError(
            "Couldn't import Django. Are you sure it's installed?"
        ) from exc
    execute_from_command_line(sys.argv)`,
		},
	}

	// GoGinTestData provides Go Gin framework test data
	GoGinTestData = TestData{
		Files: map[string]string{
			"go.mod": `module example.com/myapp

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
)`,
			"main.go": `package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, Gin!",
		})
	})
	r.Run()
}`,
		},
	}

	// LaravelTestData provides Laravel framework test data
	LaravelTestData = TestData{
		Files: map[string]string{
			"composer.json": `{
    "name": "laravel/laravel",
    "type": "project",
    "require": {
        "php": "^8.1",
        "laravel/framework": "^10.0"
    }
}`,
			"artisan": `#!/usr/bin/env php
<?php

define('LARAVEL_START', microtime(true));

require __DIR__.'/vendor/autoload.php';

$app = require_once __DIR__.'/bootstrap/app.php';

$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);

$status = $kernel->handle(
    $input = new Symfony\Component\Console\Input\ArgvInput,
    new Symfony\Component\Console\Output\ConsoleOutput
);

$kernel->terminate($input, $status);`,
		},
	}
)