const express = require('express');
const bodyParser = require('body-parser');

const { createCandidateRouter } = require('./router');

const app = express();
const port = 3000;

app.use(bodyParser.json());
app.use(createCandidateRouter());

app.get('/', (req, res) => {
  res.send('Hello World!')
})

app.use(express.static('.'))

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})
