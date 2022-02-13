<template>
  <div>
    <v-card>
      <v-data-table
          :headers="headers"
          :items="records"
      >
        <template v-slot:top>
          <v-toolbar flat>
            <v-toolbar-title>Records</v-toolbar-title>
            <v-spacer></v-spacer>
            <v-dialog v-model="dialogNew" max-width="600px">
              <template v-slot:activator="{ on, attrs }">
                <v-btn color="primary" dark class="mb-2" v-bind="attrs" v-on="on">
                  New Record
                </v-btn>
              </template>
              <v-card>
                <v-card-title>New Record</v-card-title>
                <v-card-text>
                  <v-container>
                    <v-row>
                      <v-col cols="12" sm="6" md="4">
                        <v-text-field
                            v-model="editedRecord.rr"
                            label="Resource record"
                        ></v-text-field>
                      </v-col>
                      <v-col cols="12" sm="6" md="4">
                        <v-select
                            v-model="editedRecord.dn"
                            :items="domainNames"
                            label="Domain name"
                        ></v-select>
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col cols="12" sm="6" md="4">
                        <v-text-field
                          v-model="editedRecord.value"
                          label="Record value"
                        ></v-text-field>
                      </v-col>
                      <v-col cols="12" sm="6" md="4">
                        <v-select
                            v-model="editedRecord.type"
                            :items="recordTypes"
                            label="Record type"
                        ></v-select>
                      </v-col>
                    </v-row>
                  </v-container>
                </v-card-text>
                <v-card-actions>
                  <v-spacer></v-spacer>
                  <v-btn color="blue darken-1" text @click="closeDialogNew">Cancel</v-btn>
                  <v-btn color="blue darken-1" text @click="saveRecord" :loading="savingRecord">Save</v-btn>
                </v-card-actions>
              </v-card>
            </v-dialog>
            <v-dialog v-model="dialogDelete" max-width="500px">
              <v-card>
                <v-card-title>Are you sure you want to delete this record?</v-card-title>
                <v-card-actions>
                  <v-spacer></v-spacer>
                  <v-btn color="blue darken-1" text @click="closeDialogDelete">Cancel</v-btn>
                  <v-btn color="blue darken-1" text @click="deleteRecordConfirm" :loading="deletingRecord">OK</v-btn>
                  <v-spacer></v-spacer>
                </v-card-actions>
              </v-card>
            </v-dialog>
          </v-toolbar>
        </template>
        <template v-slot:item.actions="{ item }">
          <v-icon small @click="deleteRecord(item)">
            mdi-delete
          </v-icon>
        </template>
      </v-data-table>
    </v-card>
    <v-snackbar v-model="showMessage" timeout="3000" top>{{ message }}</v-snackbar>
  </div>
</template>

<script>
import api from "@/api";

export default {
  name: "Records",

  data() {
    return {
      message: null,
      showMessage: false,
      savingRecord: false,
      deletingRecord: false,
      dialogNew: false,
      dialogDelete: false,
      headers: [
        {
          text: 'Resource Record',
          align: 'start',
          sortable: false,
          value: 'rr',
        },
        { text: 'Domain Name', value: 'dn' },
        { text: 'Type', value: 'type'},
        { text: 'Value', value: 'value'},
        { text: 'Actions', value: 'actions', sortable: false },
      ],
      domainNames: [],
      recordTypes: [],
      records: [],
      editedIndex: -1,
      editedRecord: {
        rr: '',
        dn: '',
        type: '',
        value: ''
      },
      defaultRecord: {
        rr: '',
        dn: '',
        type: '',
        value: ''
      },
    }
  },

  created() {
    this.initialize()
  },

  watch: {
    dialogNew (val) {
      val || this.closeDialogNew()
    },
    dialogDelete (val) {
      val || this.closeDialogDelete()
    },
  },

  methods: {
    initialize() {
      api.dns.listRecords().then(res => {
        this.records = res
      })

      api.dns.listDomains().then(res => {
        let i
        for (i in res) {
          this.domainNames.push(res[i].dn)
        }
      })
      this.recordTypes = ['A', 'AAAA', 'NS', 'CNAME']
    },

    closeDialogNew() {
      this.dialogNew = false
      this.$nextTick(() => {
        this.editedRecord = Object.assign({}, this.defaultRecord)
        this.editedIndex = -1
      })
    },

    closeDialogDelete() {
      this.dialogDelete = false
      this.$nextTick(() => {
        this.editedRecord = Object.assign({}, this.defaultRecord)
        this.editedIndex = -1
      })
    },

    deleteRecordConfirm() {
      this.deletingRecord = true
      api.dns.deleteRecord(this.editedRecord.dn, this.editedRecord.rr).then(() => {
        this.records.splice(this.editedIndex, 1)
      }).catch(error => {
        console.log(this.message)
        this.message = error.response.data.error
        this.showMessage = true
      }).finally(() => {
        this.deletingRecord = false
      })
      this.closeDialogDelete()
    },

    saveRecord() {
      this.savingRecord = true
      if (this.editedIndex > -1) {
        Object.assign(this.records[this.editedIndex], this.editedRecord)
      } else {
        const editedRecord = this.editedRecord
        api.dns.addRecord(editedRecord).then(() => {
          this.records.push(editedRecord)
        }).catch(error => {
          console.log(this.message)
          this.message = error.response.data.error
          this.showMessage = true
        }).finally(() => {
          this.savingRecord = false
        })
      }
      this.closeDialogNew()
    },

    deleteRecord(item) {
      this.editedIndex = this.records.indexOf(item)
      this.editedRecord = Object.assign({}, item)
      this.dialogDelete = true
    }
  }
}
</script>

<style scoped>

</style>