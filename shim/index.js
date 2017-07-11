const { spawn } = require('child_process');
const fs = require('fs');

exports.{{.Name}} = function (req, res) {
    const main = spawn('./main');
    main.stdin.write(JSON.stringify(req.body));
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
