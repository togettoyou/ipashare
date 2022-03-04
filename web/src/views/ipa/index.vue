<template>
  <div class="app-container">
    <el-form :inline="true" :model="queryInfo" class="demo-form-inline" size="mini" ref="searchRef">
      <el-form-item label="" prop="content">
        <el-input v-model="queryInfo.content" placeholder="应用名/简介" clearable></el-input>
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
      <el-table-column prop="uuid" label="应用UUID" align="center" header-align="center"></el-table-column>
      <el-table-column prop="name" label="应用名" align="center" width="70px" header-align="center"></el-table-column>
      <el-table-column
        prop="icon_url"
        label="图标"
        align="center"
        width="70px">
        <template slot-scope="scope">
          <el-avatar shape="square" :size="48" fit="fill" :src="scope.row.icon_url"></el-avatar>
        </template>
      </el-table-column>
      <el-table-column prop="bundle_identifier" label="包名" align="center" header-align="center"></el-table-column>
      <el-table-column prop="version" label="版本" align="center" header-align="center"></el-table-column>
      <el-table-column prop="build_version" label="Build版本" align="center" header-align="center"></el-table-column>
      <el-table-column prop="mini_version" label="最低支持iOS版本" align="center" header-align="center"></el-table-column>
      <el-table-column :show-overflow-tooltip='true' prop="summary" label="简介" align="center"
                       header-align="center"></el-table-column>
      <el-table-column prop="size" label="大小" align="center" header-align="center"></el-table-column>
      <el-table-column prop="count" label="下载量" align="center" header-align="center"></el-table-column>
      <el-table-column prop="CreatedAt" label="添加时间" align="center" header-align="center"></el-table-column>
      <el-table-column
        fixed="right"
        label="功能"
        align="center"
        width="83">
        <template slot-scope="scope">
          <el-button type="text" @click="qrCode(scope.row)" size="small">安装</el-button>
          <el-button type="text" @click="download(scope.row.uuid)" size="small">下载</el-button>
        </template>
      </el-table-column>
      <el-table-column
        fixed="right"
        label="操作"
        align="center"
        width="83">
        <template slot-scope="scope">
          <el-button type="text" size="small">修改</el-button>
          <el-button type="text" size="small">删除</el-button>
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

    <el-dialog title="扫码安装" :visible.sync="dialogQrcode" center>
      <div id="qrcode"></div>
      <h2 id="name"></h2>
    </el-dialog>

    <el-dialog title="上传IPA" :visible.sync="dialogFormVisible">
      <el-form ref="form" :model="form"
               :rules="formRules">
        <el-form-item style="text-align:center;">
          <el-upload
            class="upload-demo"
            action=""
            accept=".ipa"
            :multiple=false
            :auto-upload=false
            :file-list="fileList"
            :on-change="handleChange"
          >
            <i class="el-icon-upload"></i>
            <div class="el-upload__text"><em>点击上传</em></div>
          </el-upload>
        </el-form-item>
        <el-form-item label="简介" prop="summary">
          <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 6}" v-model="form.summary"
                    placeholder="请输入简介"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="clearForm">取 消</el-button>
        <el-button type="primary" @click="upload">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {download, list, upload} from "@/api/ipa";
import QRCode from 'qrcodejs2';

export default {
  name: 'IPA',
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
        summary: '',
      },
      formRules: {
        summary: [{required: true, trigger: "blur", min: 2, max: 100, message: "请输入简介（2-100字数）"}],
      },
      fileList: [],
      dialogQrcode: false,
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
            message: '请上传IPA文件',
            type: 'warning'
          });
          return false;
        }
        if (valid) {
          const formData = new FormData()
          formData.append('ipa', this.fileList[0].raw)
          formData.append('summary', this.form.summary)
          upload(formData).then(res => {
            console.log(res);
            this.$message.success('上传成功')
            this.clearForm()
          }).catch(err => {
            console.log(err);
          })
        } else {
          return false;
        }
      });
    },
    clearForm() {
      this.$refs.form.resetFields()
      this.fileList = []
      this.dialogFormVisible = false
    },
    handleChange(file, fileList) {
      if (fileList.length > 0) {
        this.fileList = [fileList[fileList.length - 1]]
      }
    },
    qrCode(row) {
      this.dialogQrcode = true;
      this.$nextTick(function () {
        document.getElementById("qrcode").innerHTML = "";
        document.getElementById("name").innerHTML = row.name + " V" + row.version;
        new QRCode("qrcode", {
          width: 150,
          height: 150,
          text: row.install_url,
        });
      });
    },
    download(uuid) {
      download(uuid).then(res => {
        console.log(res);
        let data = res
        if (!data) {
          return
        }
        let url = window.URL.createObjectURL(new Blob([data]))
        let a = document.createElement('a')
        a.style.display = 'none'
        a.href = url
        a.setAttribute('download', uuid + ".ipa")
        document.body.appendChild(a)
        a.click()
        window.URL.revokeObjectURL(a.href)
        document.body.removeChild(a)
      }).catch(err => {
        console.log(err);
      })
    }
  },
}
</script>

<style lang="scss" scoped>
</style>
