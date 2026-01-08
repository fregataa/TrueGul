# TrueGul ML Server

AI writing detection and feedback service.

## Requirements

- Python 3.13+

## Development Setup

### 1. Install Dependencies

```bash
# Using virtual environment (recommended)
python -m venv .venv
source .venv/bin/activate  # Linux/macOS
# .venv\Scripts\activate   # Windows

pip install -e ".[dev]"
```

### 2. Code Quality Tools

#### Individual Commands

```bash
# Linter (Ruff)
ruff check app/          # Check only
ruff check --fix app/    # Auto-fix

# Formatter (Ruff)
ruff format --check app/  # Check only
ruff format app/          # Auto-format

# Type Checker (Pyright)
pyright app/
```

#### Using uvx (without virtual environment)

```bash
uvx ruff check app/
uvx ruff format app/
uvx pyright app/
```

#### Full Check (same as CI)

```bash
ruff check app/ && ruff format --check app/ && pyright app/
```

### 3. Pre-commit Hooks

```bash
# Install hooks (one-time)
pre-commit install

# Run manually
pre-commit run --all-files
```

### 4. Run Tests

```bash
pytest -v
```

## Running the Server

```bash
uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload
```
