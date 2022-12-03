const express = require("express"),
      bodyParser = require("body-parser"),
      path = require("path"),
      webpush = require("web-push");

const vapidKeys = webpush.generateVAPIDKeys(),
      port = 1234,
      mail = "qwerty",
      publicKey = "BCgL7bR1TaXuKcAU4lT15tsdgG5EQg3g2KOrOmd7KtY-Z-A30lG4FaJf9OxasORYlgv-xD-KAp8rZ0aB-1bHYQU",
      privateKey = "3fD_mIAwxG-L8Lb-UnxN0KV6UQOPldxzO58qpFZUWHk";

webpush.setVapidDetails("mailto:" + mail, publicKey, privateKey);
const app = express();
app.use(express.static(__dirname));
app.use(bodyParser.json());

app.get("/", (req, res) => {
  res.sendFile(path.join(__dirname, "/index.html"));
});

app.post("/sendpush", (req, res) => {
  res.status(201).json({}); 
  webpush.sendNotification(req.body, JSON.stringify({
    title: "Notification",
    body: "Text"
  }))
  .catch((err) => { console.log(err); });
});

app.listen(port, () => {
  console.log(`Listening on ${port}`)
});
