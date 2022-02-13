function getBearerToken() {
    const token = window.localStorage.getItem('jwt')
    return token ? `Bearer ${token}` : null
}

export default {
    getBearerToken
}