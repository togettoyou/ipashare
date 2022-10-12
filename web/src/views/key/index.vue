<template>
  <div class="app-container">

    <el-switch
      style="display: block"
      v-model="enable_key"
      @change="changeKeyEnable"
      validate-event
      active-color="#13ce66"
      inactive-color="#ff4949"
      active-text="密钥功能已开启"
      inactive-text="密钥功能已关闭">
    </el-switch>
    <br><br>
    <el-form :inline="true" :model="queryInfo" class="demo-form-inline" size="mini" ref="searchRef">
      <el-form-item label="" prop="content">
        <el-input v-model="queryInfo.content" placeholder="用户名" clearable></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="clearSearch">清空</el-button>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="search">查询</el-button>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="dialogFormVisible = true">新增</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="List" style="width: 100%;" stripe highlight-current-row border
              :header-cell-style="{background:'#f8f8f9',color: '#606266','font-weight':'bold'}">
      <el-table-column prop="username" label="用户名" align="center" header-align="center"></el-table-column>
      <el-table-column prop="password" label="密码" align="center" header-align="center"></el-table-column>
      <el-table-column prop="num" label="剩余可使用个数" align="center" header-align="center"></el-table-column>
      <el-table-column prop="created_at" label="添加时间" align="center" header-align="center"></el-table-column>
      <el-table-column
        fixed="right"
        label="操作"
        align="center"
        width="180">
        <template slot-scope="scope">
          <el-button type="text" @click="updateDialog(scope.row)" size="small">修改</el-button>
          <el-button type="text" @click="del(scope.row.username)" size="small">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <br>
    <el-pagination
      background
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="queryInfo.page"
      :page-sizes="[2, 5, 10, 20, 30]"
      :page-size="queryInfo.page_size"
      layout="total, sizes, prev, pager, next, jumper"
      :total="total">
    </el-pagination>

    <el-dialog title="修改密钥可用个数" :visible.sync="updateDialogFormVisible"
               :show-close="false"
               :close-on-click-modal="false"
               :close-on-press-escape="false"
               :destroy-on-close="true"
    >
      <el-form ref="updateForm" :model="updateForm">
        <el-form-item prop="num">
          <el-input-number v-model="updateForm.num" :min="0"></el-input-number>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="clearUpdateForm">取 消</el-button>
        <el-button type="primary" @click="update" :loading="uploading">确 定</el-button>
      </div>
    </el-dialog>


    <el-dialog title="新增密钥" :visible.sync="dialogFormVisible"
               :show-close="false"
               :close-on-click-modal="false"
               :close-on-press-escape="false"
               :destroy-on-close="true"
    >
      <el-form ref="form" :model="form" id="dialogForm" label-position="right" label-width="120px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username"
                    placeholder="请输入用户名" minlength="4" maxlength="50" show-word-limit></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password"
                    placeholder="请输入密码" minlength="4" maxlength="50" show-word-limit></el-input>
        </el-form-item>
        <el-form-item label="可使用个数" prop="num">
          <el-input-number v-model="form.num" :min="0"></el-input-number>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="clearForm">取 消</el-button>
        <el-button type="primary" @click="uuidKey" :loading="uploading">使用安全级别密钥填充（建议）</el-button>
        <el-button type="primary" @click="createKey" :loading="uploading">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {addKey, delKey, keyList, updateKey} from "@/api/key";
import {getKeyConf, setKeyConf} from "@/api/conf";
import {Loading} from 'element-ui';
import {v4 as uuidv4} from 'uuid';

export default {
  name: 'Key',
  data() {
    return {
      enable_key: false,
      List: [],
      queryInfo: {
        page: 1,
        page_size: 5,
        content: ''
      },
      total: 0,
      dialogFormVisible: false,
      form: {
        username: '',
        password: '',
        num: 0,
      },
      uploading: false,
      uploadCancel: null,
      dialogQrcode: false,
      updateDialogFormVisible: false,
      updateForm: {
        username: '',
        num: 0,
      },
    }
  },
  created() {
    this.getListFilter()
  },
  methods: {
    changeKeyEnable() {
      let msg = ''
      console.log(this.enable_key)
      if (this.enable_key) {
        msg = '即将开启密钥功能, 后续扫码下载需要验证密钥, 是否继续?'
      } else {
        msg = '即将关闭密钥功能, 后续扫码下载无需验证密钥, 是否继续?'
      }
      this.$confirm(msg, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        setKeyConf({'enable_key': this.enable_key}).then(res => {
          this.$message.success("更改成功")
          this.getListFilter()
        }).catch(err => {
          this.$message.error("更改失败" + err)
          this.getListFilter()
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消更改'
        });
        this.enable_key = !this.enable_key
      });
    },
    getListFilter() {
      keyList(this.queryInfo).then(res => {
        console.log(res);
        if (res.code === 0) {
          this.List = res.data.data
          this.total = res.data.total
        } else {
          this.$message.info(res.msg)
        }
      }).catch(err => {
        console.log(err);
      })
      getKeyConf().then(res => {
        console.log(res);
        if (res.code === 0) {
          this.enable_key = res.data.enable_key
        } else {
          this.$message.info(res.msg)
        }
      }).catch(err => {
        console.log(err);
      })
    },
    handleSizeChange(pageSize) {
      this.queryInfo.page = 1
      this.queryInfo.page_size = pageSize
      this.getListFilter()
    },
    handleCurrentChange(page) {
      this.queryInfo.page = page
      this.getListFilter()
    },
    search() {
      this.queryInfo.page = 1
      this.getListFilter()
    },
    clearSearch() {
      this.$refs.searchRef.resetFields()
      this.getListFilter()
    },
    uuidKey() {
      this.form = {
        username: uuidv4(),
        password: uuidv4(),
        num: this.form.num,
      }
    },
    createKey() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.uploading = true
          let loading = Loading.service({
            target: "#dialogForm",
            fullscreen: true,
            text: '正在创建'
          })
          addKey(this.form).then(res => {
            this.$message.success('创建成功')
            loading.close()
            this.clearForm()
            this.getListFilter()
          }).catch(err => {
            loading.close()
            this.clearForm()
          })
        } else {
          return false;
        }
      });
    },
    updateDialog(row) {
      this.updateForm = {
        username: row.username,
        num: row.num,
      }
      this.updateDialogFormVisible = true
    },
    clearUpdateForm() {
      this.updateForm = {
        username: '',
        num: 0,
      }
      this.updateDialogFormVisible = false
      this.uploading = false
    },
    update() {
      this.$refs.updateForm.validate((valid) => {
        if (valid) {
          this.uploading = true
          updateKey(this.updateForm).then(res => {
            this.$message.success('修改成功')
            this.clearUpdateForm()
            this.getListFilter()
          }).catch(err => {
            this.clearUpdateForm()
          })
        } else {
          return false;
        }
      });
    },
    clearForm() {
      if (this.uploadCancel) {
        this.uploadCancel()
      }
      this.uploadCancel = null
      this.form = {
        username: '',
        password: '',
        num: 0,
      }
      this.dialogFormVisible = false
      this.uploading = false
    },
    del(username) {
      this.$confirm('此操作将删除该密钥, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        delKey(username).then(res => {
          this.$message.success("删除成功")
          this.getListFilter()
        }).catch(err => {
          this.$message.error("删除失败" + err)
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        });
      });
    }
  },
}
</script>

<style lang="scss" scoped>
</style>
