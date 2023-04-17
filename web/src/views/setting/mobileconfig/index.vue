<template>
  <div class="app-container">
    <el-alert title="可选。为空时描述文件默认显示未签名但不影响使用，配置后可显示为已签名。"/>
    <br>
    <el-form :model="ruleForm" status-icon ref="ruleForm" class="changeForm">
      <el-form-item label="域名证书（*.crt文件）" prop="server_crt_content">
        <el-input type="textarea" :rows="5" v-model="ruleForm.server_crt_content"></el-input>
      </el-form-item>
      <el-form-item label="密钥（*.key文件）" prop="server_key_content">
        <el-input type="textarea" :rows="5" v-model="ruleForm.server_key_content"></el-input>
      </el-form-item>
      <el-form-item label="根证书（root_bundle.crt文件）" prop="cert_chain_crt_content">
        <el-input type="textarea" :rows="5" v-model="ruleForm.cert_chain_crt_content"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submitForm()">保存</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import {getMobileConfig, setMobileConfig} from "@/api/conf";
import {Loading} from 'element-ui';

export default {
  name: 'Setting',
  data() {
    return {
      ruleForm: {
        server_crt_content: '',
        server_key_content: '',
        cert_chain_crt_content: ''
      },
    };
  },
  created() {
    this.getForm()
  },
  methods: {
    getForm() {
      getMobileConfig().then(res => {
        this.ruleForm = res.data;
      }).catch(err => {
      })
    },
    submitForm() {
      this.$confirm('确认保存?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        let loadingInstance = Loading.service({
          lock: true,
          text: '保存中......',
          spinner: 'el-icon-loading',
          background: 'rgba(0, 0, 0, 0.7)'
        });
        setMobileConfig(this.ruleForm).then(res => {
          this.$message.success('保存成功')
          this.getForm()
          loadingInstance.close()
        }).catch(err => {
          this.getForm()
          loadingInstance.close()
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消保存'
        });
      });
    },
  }
}
</script>

<style lang="scss" scoped>
.changeForm {
  width: 500px;
  padding-left: 10px;
}
</style>
