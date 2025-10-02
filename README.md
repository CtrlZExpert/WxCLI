# 🌦️ wxcli

A simple command-line weather application written in **Go**.  
This project is part of my personal journey to learn Go by building real, usable tools.  

---

## 🚀 Features

- Get current weather for any city  
- Clean CLI output with colors  
- Supports multiple locations  
- (Work in progress) More features on the roadmap below  

---

## 🗺️ Roadmap

- [ ] Unit Switching (°F ↔ °C, mph ↔ kph)  
  Allow users to toggle between metric and imperial units.  

- [ ] Search Autocomplete / History  
  Save previously searched locations and let users quickly re-select them.  

- [ ] Caching  
  Cache the most recent weather data so results load instantly when repeated within a short time window.  

---

## 📦 Installation

Clone the repo and build it:

```bash
git clone https://github.com/<your-username>/wxcli.git
cd wxcli
go build -o wxcli
```

Run it:

```bash
./wxcli <city>
```

---

## 📝 Example

```bash
./wxcli London
```

Output:

```
🌤️  London, UK
Temp: 17°C (62°F)
Wind: 12 kph
Condition: Partly Cloudy
```

---

## 🎯 Why I Built This

I wanted to practice Go by making a **practical CLI tool** that combines API calls, clean terminal output, and error handling.  
This project helps me improve at:  

- Writing idiomatic Go  
- Working with APIs  
- Handling CLI arguments  
- Building and structuring Go projects  

---

## 📌 Next Steps

- Expand feature set (see roadmap)  
- Improve CLI design and usability  
- Explore packaging/distribution  

---

## 🤝 Contributing

This project is mainly for learning, but suggestions and PRs are welcome.  

---

## 📜 License

MIT License.  
