# 🌊 Golang FiveM Flooder 2024

### Un outil de simulation de charge pour FiveM utilisant des proxies

**Version :** v1  
**Auteur :** [softwaretobi](https://github.com/softwaretobi)  

---

## 🚀 Description

Le **Golang FiveM Flooder** est un outil conçu pour simuler des charges de trafic sur les serveurs FiveM. Il utilise des proxies pour masquer les adresses IP, permettant ainsi d'effectuer des tests de résistance tout en restant discret. Cet outil est destiné à des fins éducatives et de recherche. Veuillez l'utiliser de manière responsable.

---

## ⚙️ Fonctionnalités

- **Proxy Rotation :** Utilise des proxies pour masquer l'adresse IP d'origine.
- **Modes de fonctionnement :** Choisissez entre un mode "stealth" (discret) et un mode "flood" (inondation).
- **Suivi des performances :** Surveille l'utilisation de la mémoire et le nombre de requêtes réussies.
- **Support multi-thread :** Effectue des requêtes simultanées pour des tests plus efficaces.

---

## 📦 Installation

1. Clonez ce dépôt :
   ```bash
   git clone https://github.com/softwaretobi/fivem-flooder.git
   cd fivem-flooder
   ```

2. Assurez-vous d'avoir Go installé sur votre machine. Si ce n'est pas le cas, téléchargez-le depuis [golang.org](https://golang.org/dl/).

3. Compilez le programme :
   ```bash
   go build -o fivem_flooder
   ```

---

## 📜 Utilisation

```bash
go run script.go [IP] [Game Port] [Mode True pour stealth, False pour flood] [Durée en secondes] [chemin_du_fichier_proxies] [nombre_de_threads] [Txadmin PORT]
```

### Paramètres :

- **IP** : L'adresse IP du serveur cible.
- **Game Port** : Le port du serveur de jeu.
- **Mode** : `true` pour le mode stealth, `false` pour le mode flood.
- **Durée** : Durée en secondes pour l'exécution de l'outil.
- **proxies_file_path** : Chemin vers un fichier contenant la liste des proxies.
- **threads** : Nombre de threads à utiliser.
- **Txadmin PORT** : Port de TxAdmin, par défaut à `Game Port + 10000`.

---

## 🔍 Exemples

Pour lancer le flooder avec les paramètres suivants :

```bash
go run script.go 192.168.1.1 30120 false 60 proxies.txt 10 30130
```

Cela va envoyer des requêtes flood au serveur avec l'IP `192.168.1.1` sur le port `30120` pendant `60` secondes en utilisant `10` threads et les proxies spécifiés dans `proxies.txt`.

---

## 🔒 Avertissement

**Cet outil est uniquement destiné à des fins éducatives.** L'utilisation de ce type de logiciel contre des serveurs sans autorisation est illégale et contraire à l'éthique. Vous êtes responsable de l'utilisation que vous en faites.

---

## 📄 Licence

Ce projet est sous licence MIT - voir le fichier [LICENSE](LICENSE) pour plus de détails.

---

## 📞 Contact

Pour toute question ou contribution, n'hésitez pas à me contacter via [GitHub](https://github.com/softwaretobi).
