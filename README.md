
# Izumi

<p align="center">
  <img src="./doc/images/pwa-512x512.png" width="28%" alt="Izumi Logo" />
</p>

<p align="center">
<strong>Izumi</strong> is a cross-platform, feature-rich game metadata scraper and management service designed to help you easily organize and enrich your local game library.
</p>

<p align="center">
  <img src="./doc/images/scrap.png" width="70%" alt="Scraper Preview" />
</p>

---

## ✨ Features

### 🧩 Built-in Metadata Scrapers
Izumi includes six powerful scrapers out of the box:

- **Bangumi**
- **DLsite**
- **Getchu**
- **2DFan**
- **GGBases**
- **VNDB**

### 📦 Rich Metadata Support
Fetch and manage:

- Tags / Categories  
- Brands / Publishers  
- Character information  
- Basic & extended metadata  
- Additional fields depending on the source website

### 🔁 Local Game Library Integration
- Automatically save scraped metadata into your game directory  
- Import existing metadata from game folders  

### 🔍 Search Functionality
Search across supported sites and your local database.

---

## 📸 Gallery

<p align="center">
  <img src="./doc/images/games.png" width="85%" alt="Games View" />
</p>

<p align="center">
  <img src="./doc/images/gamedetail1.png" width="85%" alt="Game Detail 1" />
</p>

<p align="center">
  <img src="./doc/images/gamedetail2.png" width="85%" alt="Game Detail 2" />
</p>

---

## 🚧 TODO

- [ ] Implement the Dashboard panel (currently placeholder data)
- [ ] Add English localization

---

## 🚀 Installation

### Linux & macOS

```sh
git clone git@github.com:dokidokikoi/izumi.git
cd izumi
docker-compose build
docker-compose up -d
```

Then open:
👉 http://localhost:8080