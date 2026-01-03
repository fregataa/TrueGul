import pytest
from fastapi.testclient import TestClient

from app.main import app

client = TestClient(app)


def test_health_check():
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "healthy"}


def test_root():
    response = client.get("/")
    assert response.status_code == 200
    assert response.json() == {"message": "TrueGul ML Server"}


def test_analyze_empty_text():
    response = client.post("/api/v1/analyze", json={"text": "", "writing_type": "essay"})
    assert response.status_code == 400


def test_analyze_valid_text():
    response = client.post("/api/v1/analyze", json={"text": "This is a test.", "writing_type": "essay"})
    assert response.status_code == 200
    data = response.json()
    assert "ai_score" in data
    assert "feedback" in data
    assert "model_version" in data
