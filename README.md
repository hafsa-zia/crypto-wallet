
# 🚀 **Crypto Wallet System – Custom Blockchain + UTXO Model + PoW Mining**

A fully functional **cryptocurrency wallet system** built using:

* **Backend:** Go (Golang), Gin, MongoDB, JWT, AES encryption
* **Frontend:** React + Vite + Tailwind + Context API
* **Blockchain:** Custom blockchain, UTXO model, Proof-of-Work mining
* **Security:** RSA keys, AES-encrypted private keys, email OTP verification
* **Deployment:** Backend on Railway, Frontend on Vercel

---

## 📌 **Overview**

This project implements a **real blockchain workflow** including:

* Wallet creation with **public/private keypair**
* **AES-encrypted** private key storage
* Wallet balances derived using **UTXO**
* Creating and signing **transactions**
* **Mining blocks with Proof-of-Work**
* Global blockchain structure stored in MongoDB
* **Zakat auto-deduction (2.5%)** from balances
* System & transaction logs
* Beautiful dark-themed React UI

Essentially, this is a simplified Bitcoin-style blockchain tailored for a wallet application.

---

# 🎯 **Features**

## 🔐 **Authentication**

* Email + OTP verification
* Password hashing using bcrypt
* JWT-based session authentication
* Secure login/logout flow

## 💳 **Wallet**

* RSA public/private key generation
* Wallet ID = SHA-256(public_key)
* Private key stored encrypted (AES-256-GCM)
* Editable profile (name, CNIC, non-editable email unless re-verified)

## 💰 **UTXO-Based Balance**

* Balance = sum of unspent transaction outputs
* Prevents double spending
* Same model used in Bitcoin

## 🔄 **Transactions**

* Built using UTXO inputs + outputs
* Signed with the user’s decrypted private key
* Saved as **pending** until mined
* Includes:

  * Normal transfers
  * Zakat deductions
  * Mining reward

## ⛏️ **Mining (Proof-of-Work)**

* Collects all pending transactions
* Creates a block with:

  * previous hash
  * transactions
  * timestamp
  * nonce
* Runs SHA-256 until hash starts with `"00000"`
* Saves block in chain and confirms transactions

## 📊 Reports & Explorer

* **Block Explorer** (view whole chain)
* **Transaction History**
* **System Logs**
* **Wallet Reports**:

  * Total sent
  * Total received
  * Total zakat deducted
  * Net flow

## 🕌 **Automated Zakat (2.5%)**

* Monthly deduction
* Creates a zakat transaction
* Added to blockchain
* Tracked in reports

---

# 🧱 **Architecture Diagram**

```
                 ┌──────────────────────────┐
                 │         Frontend         │
                 │ React + Vite + Tailwind  │
                 └──────────────┬───────────┘
                                │ API Calls
                                ▼
                 ┌──────────────────────────┐
                 │          Backend         │
                 │   Go + Gin + MongoDB     │
                 │                          │
                 │ - Auth + OTP Email       │
                 │ - Generate RSA Keys      │
                 │ - AES Encryption         │
                 │ - UTXO Engine            │
                 │ - Transaction Builder    │
                 │ - Digital Signatures     │
                 │ - Proof-of-Work Miner    │
                 │ - Blockchain Storage     │
                 └──────────────┬───────────┘
                                ▼
                 ┌──────────────────────────┐
                 │        MongoDB           │
                 │ users / wallets          │
                 │ utxos / txs / blocks     │
                 │ system logs / otp store  │
                 └──────────────────────────┘
```

---

# ⚙️ **Backend Installation**

### Clone the repo

```bash
git clone https://github.com/yourusername/crypto-wallet-backend.git
cd crypto-wallet-backend
```

### Install dependencies

```bash
go mod tidy
```

### Add your `.env`

Create a `.env` file:

```
MONGO_URI=your_mongo_connection
JWT_SECRET=your_jwt_secret
AES_SECRET_KEY=64hexcharacterslongkey
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASS=your_16_char_app_password
SMTP_FROM=your_email@gmail.com
```

### Run the backend

```bash
go run cmd/main.go
```

---

# 🎨 **Frontend Installation**

### Clone the frontend

```bash
git clone https://github.com/yourusername/crypto-wallet-frontend.git
cd crypto-wallet-frontend
```

### Install packages

```bash
npm install
```

### Create `.env`

```
VITE_API_URL=https://your-railway-backend.up.railway.app/api
```

### Start dev server

```bash
npm run dev
```

---

# ☁️ Deployment

## 🚂 Backend Deployment (Railway)

1. Push backend repo to GitHub
2. Open Railway → New Project → Deploy from GitHub
3. Add environment variables
4. Railway will automatically build and deploy

Backend URL example:

```
https://your-backend-production.up.railway.app/
```

---

## ▲ Frontend Deployment (Vercel)

1. Go to Vercel → New Project
2. Import repo
3. Add environment:

```
VITE_API_URL=https://your-backend-production.up.railway.app/api
```

4. Deploy 🎉

---

# 📬 Email OTP Setup

This project uses SMTP.
For Gmail:

1. Enable **2-Step Verification**
2. Under *App Passwords*, generate a 16-character password
3. Set in environment:

```
SMTP_USER=your_email@gmail.com
SMTP_PASS=your_16_char_password
```

---

# 📚 **Concepts Implemented**

* Public/Private key cryptography
* AES-256 encryption
* Digital signatures
* Proof-of-Work
* Blockchain linking (previous_hash)
* Genesis block
* UTXO model
* Double-spend prevention
* Mining reward
* Email OTP verification
* JWT Authentication
* Logs + Reports
* Zakat deduction logic

This system is a full learning implementation of real blockchain concepts.

---

# 🐛 Troubleshooting

### **Invalid Private Key**

Occurs if:

* AES_SECRET_KEY is incorrect
* Encrypted key stored wrongly
* Private key decoded incorrectly

### **No OTP Email**

Ensure:

* SMTP is configured
* Using Gmail App Password
* No 2FA issues

### **CORS Issues**

Make sure backend includes:

```go
router.Use(cors.Default())
```

---

# 🏁 Conclusion

This project implements a **complete custom blockchain system** with:

* Real cryptography
* Real UTXO model
* Real block mining
* Real transaction signing
* Real authentication + OTP
* Production-ready deployments


