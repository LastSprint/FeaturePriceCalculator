const UnlinkedJiraEpicLineCmpName = "fpc-unlinked-jira-epic-line"

var UnlinkedJiraEpicLineCmp = {
    props: ['epic'],
    methods: {
        timeAsHour: function (seconds) {
            console.log(seconds)
            if (seconds < 60.0) {
                return `${seconds}s`
            } else if (seconds < (60 * 60)) {
                return `${(seconds / 60.0).toFixed(1)}m`
            } else {
                return `${(seconds / 60.0 / 60.0).toFixed(1)}h`
            }
        },
    },
    template: `
    <vs-card style="background-color: #202020" actionable class="cardx">
        <h3 style="display: flex;justify-content: left;align-items: center;padding-left: 8px;font-family: 'Fira Code', monospace;color: #f5f3f1">
            <a :href="'https://jira.surfstudio.ru/browse/'+epic.JiraKey">{{epic.Name}}</a> 
            <p style="padding-left: 8px"> {{timeAsHour(epic.TimeSpendSum)}} </p>
        </h3>
    </vs-card>
    `
}

const UnlinkedCardCmpName = "fpc-unlinked-card"

var UnlinkedCardCmp = {
    props: ['epics'],
    methods: {
        timeAsHour: function (seconds) {
            console.log(seconds)
            if (seconds < 60.0) {
                return `${seconds}s`
            } else if (seconds < (60 * 60)) {
                return `${(seconds / 60.0).toFixed(1)}m`
            } else {
                return `${(seconds / 60.0 / 60.0).toFixed(1)}h`
            }
        },
    },
    template: `
    <vs-row vs-justify="center" style="padding-left: 22px">
        <vs-col type="flex" vs-align="left" style="width: 90%">
            <vs-card style="background-color: #414141" actionable class="cardx">
                <vs-collapse>
                    <vs-collapse-item>
                        <div slot="header">
                            <h3 style="display: flex;justify-content: left;align-items: center;padding-left: 8px;font-family: 'Fira Code', monospace;color: #f5f3f1">
                                Часов не залинковано: {{ timeAsHour(epics.reduce((a,b) => a + b.TimeSpendSum, 0)) }}
                            </h3>
                        </div>
                        <fpc-unlinked-jira-epic-line
                            v-for="it in epics"
                            v-bind:epic="it"
                            v-bind:key="'unlinked'+it.JiraKey"
                        ></fpc-unlinked-jira-epic-line>
                    </vs-collapse-item>
                </vs-collapse>
            </vs-card>
        </vs-col>
    </vs-row>
    `
}
