<template>
  <div class="app-container">
    <el-form :inline="true" :model="queryInfo" class="demo-form-inline" size="mini" ref="searchRef">
      <el-form-item label="" prop="content">
        <el-input v-model="queryInfo.content" placeholder="iss" clearable></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="clearSearch">清空</el-button>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="search">查询</el-button>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="dialogFormVisible = true">上传</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="List" style="width: 100%;" stripe highlight-current-row border
              :header-cell-style="{background:'#f8f8f9',color: '#606266','font-weight':'bold'}">
      <el-table-column prop="iss" label="iss" align="center" header-align="center"></el-table-column>
      <el-table-column prop="kid" label="kid" align="center" header-align="center"></el-table-column>
      <el-table-column prop="count" label="已使用设备量" align="center" header-align="center"></el-table-column>
      <el-table-column prop="created_at" label="添加时间" align="center" header-align="center"></el-table-column>
      <el-table-column label="已绑定设备列表" align="center" width="150">
        <template slot-scope="scope">
          <el-button type="text" @click="getDevices(scope.row.iss)" size="small">查看</el-button>
        </template>
      </el-table-column>
      <el-table-column label="最大限制量" align="center" width="100">
        <template scope="scope">
          <el-input size="small" v-model="scope.row.limit" placeholder="最大限制量"
                    type='number'
                    @change="update(scope.row)">
          </el-input>
        </template>
      </el-table-column>

      <el-table-column label="是否启用" align="center" width="100">
        <template scope="scope">
          <el-switch
            v-model="scope.row.enable"
            active-color="#13ce66"
            inactive-color="#ff4949"
            @change="update(scope.row)"
          >
          </el-switch>
        </template>
      </el-table-column>

      <el-table-column
        fixed="right"
        label="操作"
        align="center"
        width="83">
        <template slot-scope="scope">
          <el-button type="text" @click="del(scope.row.iss)" size="small">删除</el-button>
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

    <el-dialog :title="'iss['+devicesIss+'] 已绑定设备列表'" :visible.sync="devicesDialog"
               :show-close="false"
               :close-on-click-modal="false"
               :close-on-press-escape="false"
               :destroy-on-close="true"
               width="80%"
    >
      <el-table :data="devicesList" height="350" style="width: 100%;" stripe highlight-current-row border
                :header-cell-style="{background:'#f8f8f9',color: '#606266','font-weight':'bold'}">
        <el-table-column width="50" type="index" label="序号" align="center" header-align="center"></el-table-column>
        <el-table-column prop="device_id" label="设备 ID" align="center" header-align="center"></el-table-column>
        <el-table-column prop="udid" label="设备 UDID" align="center" header-align="center"></el-table-column>
        <el-table-column prop="addedDate" label="添加时间" align="center" header-align="center"></el-table-column>
        <el-table-column prop="name" label="名称" align="center" header-align="center"></el-table-column>
        <el-table-column prop="deviceClass" label="类型" align="center" header-align="center"></el-table-column>
        <el-table-column prop="deviceModel" label="型号" align="center" header-align="center"></el-table-column>
        <el-table-column prop="platform" label="平台" align="center" header-align="center"></el-table-column>
        <el-table-column prop="status" label="状态" align="center" header-align="center"></el-table-column>
      </el-table>

      <div slot="footer" class="dialog-footer">
        <el-button @click="clearDevices">关 闭</el-button>
        <el-button type="primary" :loading="devicesUploading" @click="updateDevices">从苹果后台同步最新设备</el-button>
      </div>
    </el-dialog>

    <el-dialog title="上传开发者账号" :visible.sync="dialogFormVisible"
               :show-close="false"
               :close-on-click-modal="false"
               :close-on-press-escape="false"
               :destroy-on-close="true"
    >
      <el-form ref="form" :model="form" id="dialogForm"
               :rules="formRules">
        <el-form-item style="text-align:center;">
          <el-upload
            class="upload-demo"
            action=""
            accept=".p8"
            drag
            :multiple=false
            :auto-upload=false
            :file-list="fileList"
            :on-change="handleChange"
            :on-remove="handleRemove"
          >
            <i class="el-icon-upload"></i>
            <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
          </el-upload>
        </el-form-item>
        <el-form-item label="iss" prop="iss">
          <el-input v-model="form.iss" placeholder="请输入iss"></el-input>
        </el-form-item>
        <el-form-item label="kid" prop="kid">
          <el-input v-model="form.kid" placeholder="请输入kid"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="clearForm">取 消</el-button>
        <el-button type="primary" @click="upload" :loading="uploading">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {del, list, update, upload} from "@/api/developer";
import {deviceList, updateDevice} from "@/api/device";
import {Loading} from 'element-ui';
import axios from "axios";

export default {
  name: 'Developer',
  data() {
    return {
      List: [],
      queryInfo: {
        page: 1,
        page_size: 5,
        content: ''
      },
      total: 0,
      dialogFormVisible: false,
      form: {
        iss: '',
        kid: '',
      },
      formRules: {
        iss: [{required: true, trigger: "blur", message: "请输入iss"}],
        kid: [{required: true, trigger: "blur", message: "请输入kid"}],
      },
      fileList: [],
      uploading: false,
      uploadCancel: null,
      dialogQrcode: false,
      devicesIss: '',
      devicesDialog: false,
      devicesUploading: false,
      devicesList: [],
    }
  },
  created() {
    this.getListFilter()
  },
  methods: {
    getListFilter() {
      list(this.queryInfo).then(res => {
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
    upload() {
      this.$refs.form.validate((valid) => {
        if (this.fileList.length === 0) {
          this.$message({
            message: '请上传P8文件',
            type: 'warning'
          });
          return false;
        }
        if (valid) {
          this.uploading = true
          let that = this
          let loading = Loading.service({
            target: "#dialogForm",
            fullscreen: true,
            text: '正在上传'
          })
          const formData = new FormData()
          formData.append('p8', this.fileList[0].raw)
          formData.append('iss', this.form.iss)
          formData.append('kid', this.form.kid)
          upload(formData, progressEvent => {
            loading.setText(`已上传 ${((progressEvent.loaded / progressEvent.total) * 100).toFixed(2)}%`)
            if (progressEvent.loaded === progressEvent.total) {
              loading.setText("正在校验账号可用性")
            }
          }, new axios.CancelToken(function executor(c) {
            that.uploadCancel = c;
          })).then(res => {
            this.$message.success('上传成功')
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
    clearForm() {
      if (this.uploadCancel) {
        this.uploadCancel()
      }
      this.uploadCancel = null
      this.$refs.form.resetFields()
      this.fileList = []
      this.dialogFormVisible = false
      this.uploading = false
    },
    handleChange(file, fileList) {
      if (fileList.length > 0) {
        this.fileList = [fileList[fileList.length - 1]]
      }
    },
    handleRemove() {
      this.fileList = []
    },
    getDevices(iss) {
      this.devicesIss = iss
      this.devicesDialog = true
      deviceList(iss).then(res => {
        console.log(res);
        if (res.code === 0) {
          this.devicesList = res.data
        } else {
          this.$message.info(res.msg)
        }
      }).catch(err => {
        console.log(err);
      })
    },
    clearDevices() {
      this.devicesIss = ''
      this.devicesDialog = false
      this.devicesList = []
      this.devicesUploading = false
    },
    updateDevices() {
      if (this.devicesIss !== '') {
        this.$confirm('此操作将从苹果后台同步最新的设备列表, 是否继续?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          this.devicesUploading = true
          let loading = Loading.service({
            target: "#dialogForm",
            fullscreen: true,
            text: '正在同步'
          })
          updateDevice(this.devicesIss).then(res => {
            console.log(res);
            if (res.code === 0) {
              this.$message.success("同步成功")
              loading.close()
              this.clearDevices()
            } else {
              this.$message.error("同步失败" + res.msg)
              loading.close()
              this.devicesUploading = false
            }
          }).catch(err => {
            this.$message.error("同步失败" + err)
            loading.close()
            this.devicesUploading = false
          })
        }).catch(() => {
          this.$message({
            type: 'info',
            message: '已取消同步'
          });
        });
      }
    },
    del(iss) {
      this.$confirm('此操作将永久删除该账号, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        del(iss).then(res => {
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
    },
    update(row) {
      this.$confirm('是否更改账号设置?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        update(row.iss, row.limit, row.enable).then(res => {
          this.$message.success("更改成功")
          this.getListFilter()
        }).catch(err => {
          this.$message.error("更改失败" + err)
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消'
        });
        this.getListFilter()
      });
    }
  },
}
</script>

<style lang="scss" scoped>
</style>
