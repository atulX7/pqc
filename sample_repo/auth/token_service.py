from cryptography.hazmat.primitives.asymmetric import rsa

JWT_ALGORITHM = "RS256"

def issue_token(payload):
    return {"payload": payload, "alg": JWT_ALGORITHM}
