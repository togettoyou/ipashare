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
    </el-form>
    <el-table :data="List" style="width: 100%;" stripe highlight-current-row border
              :header-cell-style="{background:'#f8f8f9',color: '#606266','font-weight':'bold'}">
      <el-table-column prop="uuid" label="应用UUID" align="center" header-align="center" fixed></el-table-column>
      <el-table-column prop="name" label="应用名" align="center" width="70px" header-align="center"
                       fixed></el-table-column>
      <el-table-column prop="bundle_identifier" label="包名" align="center" width="110px" header-align="center"
                       fixed></el-table-column>
      <el-table-column prop="version" label="版本" align="center" header-align="center"></el-table-column>
      <el-table-column prop="build_version" label="Build版本" align="center" header-align="center"></el-table-column>
      <el-table-column prop="mini_version" label="最低支持iOS版本" align="center"
                       header-align="center"></el-table-column>
      <el-table-column prop="summary" label="简介" align="center" header-align="center"></el-table-column>
      <el-table-column prop="size" label="大小" align="center" header-align="center"></el-table-column>
      <el-table-column prop="count" label="下载量" align="center" header-align="center"></el-table-column>
      <el-table-column prop="CreatedAt" label="添加时间" align="center"
                       header-align="center"></el-table-column>
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
  </div>
</template>

<script>
import {list} from "@/api/ipa";

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
  },
}
</script>

<style lang="scss" scoped>
</style>
