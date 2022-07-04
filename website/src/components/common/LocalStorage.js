/**
 * 此模块是用于local数据存储管理的工具模块
 */
import store from 'store'

const USER_KEY = 'tirelease_user'
const USER_LOGIN = false

const storage = {

    /**
     * 存储user
     */
    saveUser(user) {
        store.set(USER_KEY, user)
        store.set(USER_LOGIN, true)
    },

    /**
     * 获取user
     */
    getUser() {
        // return JSON.parse(localStorage.getItem(USER_KEY) || '{}')
        return store.get(USER_KEY) || undefined
    }
    ,

    /**
     * Get user login status
    */
    getHasLogin() {
        return store.get(USER_LOGIN) == undefined ? false : store.get(USER_LOGIN)
    }
    ,

    /**
     * 删除user
     */
    removeUser() {
        // localStorage.removeItem(USER_KEY)
        store.remove(USER_KEY)
        store.remove(USER_LOGIN)
    }
}

export default storage
