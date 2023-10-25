const express = require('express');
const bodyParser = require('body-parser');

const app = express();
const port = 8080;
app.use(bodyParser.json());

app.post('/greeting', (req, res) => {
    console.log(req.body);
    res.status(200).send();
});

app.listen(port, () => {console.log(`Receiver is running on port ${port}`)});