module.exports = {
    configureWebpack: {
        devServer: {
            proxy: {
                '^/api': {
                    target: 'http://localhost:8090'
                },
                '^/shell': {
                    target: 'http://localhost:8091'
                }
            }
        }
    },

    transpileDependencies: [
      'vuetify'
    ]
}
