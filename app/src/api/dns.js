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
    listDomainProviders() {
      return service.get('/domain/provider/list').then(res => {
          return res.data
      })
    },

    listDomains() {
        return service.get('/domain/list').then(res => {
            return res.data
        })
    },

    listRecords() {
        return service.get(`/record/list`).then(res => {
            return res.data
        })
    },

    addDomain(domain) {
        return service.post('/domain', domain).then(res => {
            return res.data
        })
    },

    addRecord(record) {
        return service.post('/record', record).then(res => {
            return res.data
        })
    },

    deleteDomain(domainName) {
        return service.delete(`/domain?dn=${domainName}`).then(res => {
            return res.data
        })
    },

    deleteRecord(domainName, resourceRecord) {
        return service.delete(`/record?dn=${domainName}&rr=${resourceRecord}`).then(res => {
            return res.data
        })
    }
}