#https://stackoverflow.com/questions/13341955/how-to-pass-a-variable-in-a-curl-command-in-shell-scripting
curl -m 70 -X POST "${URL}" \
-H "Authorization: bearer $(gcloud auth print-identity-token)" \
-H "Content-Type: application/json" \
-d '{
      "new_banks": [
        {
          "name": "HAPPIEST",
          "new_deposits": [
            {
              "account": "1234",
              "account_type": "Checking",
              "apy": 0,
              "years": 1,
              "amount": 100
            },
            {
              "account": "1256",
              "account_type": "CD",
              "apy": 24,
              "years": 2,
              "amount": 7700
            },
            {
              "account": "1111",
              "account_type": "CD",
              "apy": 1.01,
              "years": 10,
              "amount": 27000
            }
          ]
        },
        {
          "name": "NICE",
          "new_deposits": [
            {
              "account": "1234",
              "account_type": "Brokered CD",
              "apy": 2.4,
              "years": 7,
              "amount": 10990
            }
          ]
        },
        {
          "name": "COOL",
          "new_deposits": [
            {
              "account": "1234",
              "account_type": "Brokered CD",
              "apy": 5,
              "years": 7,
              "amount": 10990
            },
            {
              "account": "9898",
              "account_type": "CD",
              "apy": 2.22,
              "years": 1,
              "amount": 5500
            }
          ]
        }
      ]
    }'