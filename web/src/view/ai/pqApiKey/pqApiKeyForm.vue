
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="id字段:" prop="id">
    <el-input v-model.number="formData.id" :clearable="true" placeholder="请输入id字段" />
</el-form-item>
        <el-form-item label="updatedAt字段:" prop="updatedAt">
    <el-date-picker v-model="formData.updatedAt" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
</el-form-item>
        <el-form-item label="deletedAt字段:" prop="deletedAt">
    <el-date-picker v-model="formData.deletedAt" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
</el-form-item>
        <el-form-item label="createdAt字段:" prop="createdAt">
    <el-date-picker v-model="formData.createdAt" type="date" style="width:100%" placeholder="选择日期" :clearable="true" />
</el-form-item>
        <el-form-item label="用户id:" prop="userId">
    <el-input v-model.number="formData.userId" :clearable="true" placeholder="请输入用户id" />
</el-form-item>
        <el-form-item label="ai模型:" prop="aiModelId">
    <el-input v-model.number="formData.aiModelId" :clearable="true" placeholder="请输入ai模型" />
</el-form-item>
        <el-form-item label="密钥:" prop="key">
    <el-input v-model="formData.key" :clearable="true" placeholder="请输入密钥" />
</el-form-item>
        <el-form-item label="拥有tokens数:" prop="totalTokens">
    <el-input v-model.number="formData.totalTokens" :clearable="true" placeholder="请输入拥有tokens数" />
</el-form-item>
        <el-form-item label="已消耗tokens:" prop="useTokens">
    <el-input v-model="formData.useTokens" :clearable="true" placeholder="请输入已消耗tokens" />
</el-form-item>
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
  createPqApiKey,
  updatePqApiKey,
  findPqApiKey
} from '@/api/ai/pqApiKey'

defineOptions({
    name: 'PqApiKeyForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const formData = ref({
            id: undefined,
            updatedAt: new Date(),
            deletedAt: new Date(),
            createdAt: new Date(),
            userId: undefined,
            aiModelId: undefined,
            key: '',
            totalTokens: undefined,
            useTokens: '',
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findPqApiKey({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      btnLoading.value = true
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return btnLoading.value = false
            let res
           switch (type.value) {
             case 'create':
               res = await createPqApiKey(formData.value)
               break
             case 'update':
               res = await updatePqApiKey(formData.value)
               break
             default:
               res = await createPqApiKey(formData.value)
               break
           }
           btnLoading.value = false
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
