import express from "express";
import path from "path";
import { fileURLToPath } from "url";

const app = express();
const __dirname = path.dirname(fileURLToPath(import.meta.url));

// Set folder "html" sebagai static files
app.use(express.static(__dirname));

app.get("/", (req, res) => {
    res.sendFile(path.join(__dirname, "index.html"));
});

// Jalankan server di port 8080
const PORT = 8080;
app.listen(PORT, () => {
    console.log(`Server berjalan di PORT: ${PORT}`);
});
