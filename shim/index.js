const { spawn } = require('child_process');

exports.{{.Name}} = function (req, res) {
    const main = spawn('./main');

    reqGo = {};
    reqGo.baseUrl = req.baseUrl;
    reqGo.body = req.body;
    reqGo.hostname = req.hostname;
    reqGo.ip = req.ip;
    reqGo.method = req.method;
    reqGo.originalUrl = req.originalUrl;
    reqGo.query = req.query;

    main.stdin.write(JSON.stringify(reqGo));
    main.stdin.end();

    let result = '';

    main.stdout.on('data', (data) => {
        result += data;
    });

    main.stderr.on('data', (data) => {
        res.status(500).send(data)
    });

    main.on('close', (code) => {
        try {
            res.send(JSON.parse(result));
        } catch(ex) {
            res.send(result);
        }
    });
};
