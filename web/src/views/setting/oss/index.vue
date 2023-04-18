<template>
  <div class="app-container">
    <el-alert
      title="保存后请验证是否可用"
      type="warning">
    </el-alert>
    <br>
    <el-form :model="ruleForm" status-icon :rules="rules" ref="ruleForm" class="changeForm">
      <el-form-item label="启用对象存储 OSS" prop="enable_oss">
        <el-switch
          v-model="ruleForm.enable_oss"
          active-color="#13ce66"
          inactive-color="#ff4949"
        >
        </el-switch>
      </el-form-item>
      <el-form-item label="存储空间（Bucket）" prop="oss_bucket_name">
        <el-input v-model="ruleForm.oss_bucket_name"></el-input>
      </el-form-item>
      <el-form-item label="外网访问域名（Endpoint）" prop="oss_endpoint">
        <el-input v-model="ruleForm.oss_endpoint"></el-input>
      </el-form-item>
      <el-form-item label="内网上传域名（Endpoint）（可选）" prop="oss_lan_endpoint">
        <el-input v-model="ruleForm.oss_lan_endpoint"></el-input>
      </el-form-item>
      <el-form-item label="访问密钥（AccessKeyId）" prop="oss_access_key_id">
        <el-input v-model="ruleForm.oss_access_key_id"></el-input>
      </el-form-item>
      <el-form-item label="访问密钥（AccessKeySecret）" prop="oss_access_key_secret">
        <el-input v-model="ruleForm.oss_access_key_secret"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submitForm('ruleForm')">保存</el-button>
        <el-button type="primary" @click="verify">验证</el-button>
        <el-button @click="resetForm('ruleForm')">重置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import {get, set, verify} from "@/api/conf";
import {Loading} from 'element-ui';

export default {
  name: 'Setting',
  data() {
    var validateOSS = (rule, value, callback) => {
      if (this.ruleForm.enable_oss === true && value === '') {
        callback(new Error('请填写完整'));
      } else {
        callback();
      }
    };
    return {
      ruleForm: {
        enable_oss: false,
        oss_bucket_name: '',
        oss_endpoint: '',
        oss_lan_endpoint: '',
        oss_access_key_id: '',
        oss_access_key_secret: '',
      },
      rules: {
        oss_bucket_name: {required: true, validator: validateOSS, trigger: 'blur'},
        oss_endpoint: {required: true, validator: validateOSS, trigger: 'blur'},
        oss_access_key_id: {required: true, validator: validateOSS, trigger: 'blur'},
        oss_access_key_secret: {required: true, validator: validateOSS, trigger: 'blur'},
      },
    };
  },
  created() {
    this.getForm()
  },
  methods: {
    getForm() {
      get().then(res => {
        this.ruleForm = res.data;
      }).catch(err => {
      })
    },
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$confirm('确认更改?', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            let loadingInstance = Loading.service({
              lock: true,
              text: '更改中......',
              spinner: 'el-icon-loading',
              background: 'rgba(0, 0, 0, 0.7)'
            });
            set(this.ruleForm).then(res => {
              this.$message.success('更改成功')
              this.getForm()
              loadingInstance.close()
              this.verify()
            }).catch(err => {
              this.getForm()
              loadingInstance.close()
            })
          }).catch(() => {
            this.$message({
              type: 'info',
              message: '已取消更改'
            });
          });
        }
      });
    },
    verify() {
      let loadingInstance = Loading.service({
        lock: true,
        text: '正在验证中......',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      });
      verify().then(res => {
        this.$message.success('验证成功')
        loadingInstance.close()
      }).catch(err => {
        loadingInstance.close()
      })
    },
    resetForm(formName) {
      this.$refs[formName].resetFields();
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
