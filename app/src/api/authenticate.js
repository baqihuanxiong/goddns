import axios from "axios";

const https = require('https');
const service = axios.create({
    baseURL: `${process.env.VUE_APP_API_BASE_URL}/`,
    timeout: 15000,
    httpsAgent: new https.Agent({
        rejectUnauthorized: false
    }),
});

export default {
    login(username, password) {
        return service.post('/login', {
            username, password
        }).then(res => {
            return res.data
        })
    }
}