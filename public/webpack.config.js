var config = {
    entry: ['./js/app.jsx'],
    output: {
        filename: 'bundle.js'
    },
    module: {
        loaders: [{
            test: [/\.jsx?$/],
            loaders: ['babel']
        }, {
            test: [/\.css$/, /\.scss$/],
            loader: 'style!css!sass'
        }]
    }
};

module.exports = config;
