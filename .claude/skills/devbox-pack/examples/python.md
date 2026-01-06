# Python Examples

Complete examples of Python project analysis and deployment configuration.

## Example 1: Django Application with PostgreSQL

**User request**: "Pack my Django project"

**Detected files**:
- `requirements.txt` with `Django==4.2` and `psycopg2`
- `manage.py`
- `.python-version` → 3.11

**Output**:
```json
[
  {
    "language": "python",
    "version": "3.11",
    "apt": [],
    "dev": {
      "environment": {
        "DEBUG": "True"
      },
      "setup": ["pip install -r requirements.txt", "python manage.py migrate"],
      "commands": ["python manage.py runserver 0.0.0.0:8000"]
    },
    "prod": {
      "environment": {
        "PYTHONENV": "production",
        "DEBUG": "False"
      },
      "setup": [
        "pip install -r requirements.txt",
        "python manage.py migrate",
        "python manage.py collectstatic --noinput"
      ],
      "commands": ["gunicorn --bind 0.0.0.0:8000 --workers 4 myproject.wsgi:application"]
    },
    "port": 8000,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export PYTHONENV=production\n    export DEBUG=False\n    pip install -r requirements.txt\n    python manage.py migrate\n    python manage.py collectstatic --noinput\n    gunicorn --bind 0.0.0.0:8000 --workers 4 myproject.wsgi:application\nelse\n    export DEBUG=True\n    pip install -r requirements.txt\n    python manage.py migrate\n    python manage.py runserver 0.0.0.0:8000\nfi",
    "evidence": {
      "files": ["requirements.txt", "manage.py", ".python-version"],
      "reason": "Python Django 4.2 project detected. psycopg2 uses prebuilt binaries on Debian/glibc, no system dependencies needed."
    }
  }
]
```

**Key Points**:
- APT packages: Empty array even with `psycopg2`
  - `psycopg2-binary` wheels include libpq
  - On Debian/glibc, no postgresql-dev needed
- Django-specific commands: migrate, collectstatic
- Production: gunicorn instead of dev server
- Debug mode via environment variables

---

## Example 2: FastAPI Application

**User request**: "Analyze FastAPI project"

**Detected files**:
- `requirements.txt` with `fastapi` and `uvicorn`
- `main.py`
- `pyproject.toml`

**Output**:
```json
[
  {
    "language": "python",
    "version": "3.11",
    "apt": [],
    "dev": {
      "environment": {
        "PYTHONENV": "development"
      },
      "setup": ["pip install -r requirements.txt"],
      "commands": ["uvicorn main:app --host 0.0.0.0 --port 8000 --reload"]
    },
    "prod": {
      "environment": {
        "PYTHONENV": "production"
      },
      "setup": ["pip install -r requirements.txt"],
      "commands": ["uvicorn main:app --host 0.0.0.0 --port 8000 --workers 4"]
    },
    "port": 8000,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export PYTHONENV=production\n    pip install -r requirements.txt\n    uvicorn main:app --host 0.0.0.0 --port 8000 --workers 4\nelse\n    export PYTHONENV=development\n    pip install -r requirements.txt\n    uvicorn main:app --host 0.0.0.0 --port 8000 --reload\nfi",
    "evidence": {
      "files": ["requirements.txt", "main.py", "pyproject.toml"],
      "reason": "FastAPI application detected with uvicorn ASGI server. Using default port 8000."
    }
  }
]
```

**Key Points**:
- FastAPI uses uvicorn ASGI server
- Dev: `--reload` flag for hot reload
- Prod: `--workers 4` for multiprocessing
- Port binding: `--host 0.0.0.0 --port 8000`

---

## Example 3: Flask Application

**User request**: "Pack Flask API"

**Detected files**:
- `requirements.txt` with `Flask`
- `app.py`
- `.python-version` → 3.10

**Output**:
```json
[
  {
    "language": "python",
    "version": "3.10",
    "apt": [],
    "dev": {
      "environment": {
        "FLASK_APP": "app.py",
        "FLASK_ENV": "development",
        "FLASK_DEBUG": "1"
      },
      "setup": ["pip install -r requirements.txt"],
      "commands": ["flask run --host=0.0.0.0 --port=5000"]
    },
    "prod": {
      "environment": {
        "FLASK_APP": "app.py",
        "FLASK_ENV": "production"
      },
      "setup": ["pip install -r requirements.txt"],
      "commands": ["gunicorn --bind 0.0.0.0:5000 --workers 4 app:app"]
    },
    "port": 5000,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export FLASK_APP=app.py\n    export FLASK_ENV=production\n    pip install -r requirements.txt\n    gunicorn --bind 0.0.0.0:5000 --workers 4 app:app\nelse\n    export FLASK_APP=app.py\n    export FLASK_ENV=development\n    export FLASK_DEBUG=1\n    pip install -r requirements.txt\n    flask run --host=0.0.0.0 --port=5000\nfi",
    "evidence": {
      "files": ["requirements.txt", "app.py", ".python-version"],
      "reason": "Flask application detected. Using gunicorn for production. Default Flask port is 5000."
    }
  }
]
```

**Key Points**:
- Flask default port: 5000
- Environment variables: `FLASK_APP`, `FLASK_ENV`, `FLASK_DEBUG`
- Dev: `flask run` with debug mode
- Prod: gunicorn WSGI server

---

## Python Detection Patterns

### Framework Detection
| Framework | Detection | Default Port |
|-----------|-----------|--------------|
| Django | `manage.py`, `Django` in requirements | 8000 |
| FastAPI | `fastapi` + `uvicorn` in requirements | 8000 |
| Flask | `Flask` in requirements, `app.py` | 5000 |
| Tornado | `tornado` in requirements | 8888 |
| Sanic | `sanic` in requirements | 8000 |

### Dependency Files Priority
1. `requirements.txt` (most common)
2. `pyproject.toml` (Poetry, modern projects)
3. `setup.py` (legacy)
4. `Pipfile` (Pipenv)

### Host Binding Syntax
- Django: `runserver 0.0.0.0:PORT`
- Flask: `--host=0.0.0.0 --port=PORT`
- FastAPI/Uvicorn: `--host 0.0.0.0 --port PORT`
- Gunicorn: `--bind 0.0.0.0:PORT`

### Version Detection
1. `.python-version` file (pyenv)
2. `pyproject.toml` → `requires-python` field
3. `runtime.txt` (Heroku-style)
4. Omit if uncertain

### Production WSGI/ASGI Servers
- **Django**: gunicorn (WSGI)
- **Flask**: gunicorn (WSGI)
- **FastAPI**: uvicorn (ASGI)
- **General ASGI**: uvicorn, hypercorn, daphne

### Common Environment Variables
- `PYTHONENV`: `development` or `production`
- `DEBUG`: `True` or `False`
- `DJANGO_SETTINGS_MODULE`: Django settings module
- `FLASK_APP`: Flask application entry point
- `FLASK_ENV`: Flask environment mode

### APT Dependencies
**Default: `[]` (empty array)**

Reasons:
- Python base images include common C libraries
- `Pillow`: Binary wheels for glibc (no system libs needed)
- `psycopg2-binary`: Includes libpq
- `lxml`: Binary wheels with bundled libs
- `cryptography`: Prebuilt wheels

**Only include APT packages if**:
- CLI tools needed: `postgresql-client`, `redis-tools`
- Explicitly required in docs: `ffmpeg`, specialized libs
- Using source packages instead of binary wheels

### Django-Specific Commands
**Development**:
- `python manage.py migrate`: Database migrations
- `python manage.py runserver 0.0.0.0:8000`: Dev server

**Production**:
- `python manage.py collectstatic --noinput`: Collect static files
- `gunicorn --bind 0.0.0.0:8000 project.wsgi:application`: Production server

### Package Installation Strategies
**Development**:
```bash
pip install -r requirements.txt
```

**Production**:
```bash
pip install -r requirements.txt --no-cache-dir
# or
pip install --no-cache-dir -e .  # if using setup.py/pyproject.toml
```

**Poetry**:
```bash
poetry install --no-dev  # Production (exclude dev dependencies)
poetry install  # Development
```
