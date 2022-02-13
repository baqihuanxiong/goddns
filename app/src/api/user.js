import axios from "axios";
import envService from "@/services/env";

axios.defaults.headers.common['Authorization'] = envService.getBearerToken();

const https = require('https');
const service = axios.create({
    baseURL: `${process.env.VUE_APP_API_BASE_URL}/api/v1`,
    timeout: 15000,
    httpsAgent: new https.Agent({
        rejectUnauthorized: false
    }),
});

export default {
    getCurrentUser() {
        return service.get('/user/current').then(res => {
            return res.data
        })
    }
}