import Vue from 'vue';
import Dialog from './Dialog';
import vuetify from '@/plugins/vuetify';

const Constructor = Vue.extend(Object.assign({ vuetify }, Dialog))
const Instance = new Constructor();
Instance.$mount()
document.body.appendChild(Instance.$el)

export default function dialog({
                                   title,
                                   content,
                                   cancelText,
                                   confirmText,
                                   cancel = () => {},
                                   confirm = () => {}
                               }) {
    Instance.isDialogVisible = true
    Instance.title = title
    Instance.content = content
    Instance.cancelText = cancelText
    Instance.confirmText = confirmText
    Instance.handleCancel = cancel
    Instance.handleConfirm = confirm
}