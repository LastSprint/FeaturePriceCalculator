const TableCmpName = 'fpc-table'

var TableCmp = {
    props: ['features', 'allUnUsedTime', 'timeArr'],
    methods: {
        timeAsHour: function (seconds) {
            if (seconds < 60.0) {
                return `${seconds}s`
            } else if (seconds < (60 * 60)) {
                return `${(seconds / 60.0).toFixed(1)}m`
            } else {
                return `${(seconds / 60.0 / 60.0).toFixed(1)}h`
            }
        },
        getUnused: function (f, u, i) {
            var cAll = u / f.length
            console.log("** cAll", cAll)

            var currentSpent = f[i].SpentSum

            var allSpent = f.reduce((a,b)=>a + b.SpentSum, 0)

            var sSumR = currentSpent / allSpent

            console.log("** currentSpent", currentSpent)
            console.log("** allSpent", allSpent)
            console.log("** sSumR", sSumR)
            console.log("** Res", sSumR * u)

            return sSumR * u
        },
        getDiff: function (arr, i) {
            var val = arr[i].Estimate - arr[i].SpentSum

            var time = this.timeAsHour(Math.abs(val))

            if (val < 0) {
                return "-"+ time
            }
            return time
        }
    },//202020
    template:` 
    <vs-table :data="features" style="background-color: #414141">
      <template slot="header">
        <vs-row>
            <vs-card style="background-color: #414141" actionable class="cardx">
                <h1 style="display: flex;justify-content: left;align-items: center;font-family: 'Fira Code', monospace;color: #f5f3f1;">
                    Сводная таблица
                </h1>
            </vs-card>
            <vs-card style="background-color: #202020" actionable class="cardx">
                <h3 style="display: flex;justify-content: left;align-items: center;font-family: 'Fira Code', monospace;color: #f5f3f1;">
                    Всего оценено (смета): {{ timeAsHour(features.reduce((a,b)=>a + b.Estimate, 0)) }}
                </h3>
                </br>
                <h3 style="display: flex;justify-content: left;align-items: center;font-family: 'Fira Code', monospace;color: #f5f3f1;">
                    Всего затрекано в залинкованные эпики: {{ timeAsHour(features.reduce((a,b)=>a + b.SpentSum, 0)) }}
                </h3>
                <h3 style="display: flex;justify-content: left;align-items: center;font-family: 'Fira Code', monospace;color: #f5f3f1;padding-top: 16px">
                    Не распределено: {{ timeAsHour(allUnUsedTime) }}
                </h3>
            </vs-card>
        </vs-row>
      </template>
      <template slot="thead" style="background-color: #414141">
        <vs-th style="background-color: #414141">
          <p style="display: flex;justify-content: left;align-items: center;padding: 8px;font-family: 'Fira Code', monospace;color: #f5f3f1;padding-top: 8px">
            Фича из сметы
          </p>
        </vs-th>
        <vs-th style="background-color: #414141">
          <p style="display: flex;justify-content: left;align-items: center;padding: 8px;font-family: 'Fira Code', monospace;color: #f5f3f1;padding-top: 8px">
            Оценка (Смета)
          </p>
        </vs-th>
        <vs-th style="background-color: #414141">
          <p style="display: flex;justify-content: left;align-items: center;padding: 8px;font-family: 'Fira Code', monospace;color: #f5f3f1;padding-top: 8px">
            Траты (Jira)
          </p>
        </vs-th>
        <vs-th style="background-color: #414141">
          <p style="display: flex;justify-content: left;align-items: center;padding: 8px;font-family: 'Fira Code', monospace;color: #f5f3f1;padding-top: 8px">
            Сколько еще можно затрекать
          </p>
        </vs-th>
        <vs-th style="background-color: #414141">
          <p style="display: flex;justify-content: left;align-items: center;padding: 8px;font-family: 'Fira Code', monospace;color: #f5f3f1;padding-top: 8px">
            Долитые часы (автоматически вычислено)
          </p>
        </vs-th>
        <vs-th style="background-color: #414141">
          <p style="display: flex;justify-content: left;align-items: center;padding: 8px;font-family: 'Fira Code', monospace;color: #f5f3f1;padding-top: 8px">
            Результат
          </p>
        </vs-th>
      </template>

      <template slot-scope="{data}">
        <vs-tr :key="indextr" v-for="(tr, indextr) in data" >
          <vs-td :data="data[indextr].Name" style="background-color: #414141">
            <p style="display: flex;justify-content: left;align-items: center;padding: 8px;font-family: 'Fira Code', monospace;color: #f5f3f1;padding-top: 8px">
                {{data[indextr].Name}}
            </p>
          </vs-td>

          <vs-td :data="data[indextr].Estimate" style="background-color: #414141">
            <p style="display: flex;justify-content: left;align-items: center;padding: 8px;font-family: 'Fira Code', monospace;color: #f5f3f1;padding-top: 8px">
                {{ timeAsHour(data[indextr].Estimate) }}
            </p>
          </vs-td>

          <vs-td :data="data[indextr]" style="background-color: #414141">
            <p style="display: flex;justify-content: left;align-items: center;padding: 8px;font-family: 'Fira Code', monospace;color: #f5f3f1;padding-top: 8px">
                {{ timeAsHour(data[indextr].SpentSum) }}
            </p>
          </vs-td>
          
          <vs-td :data="data" style="background-color: #414141">
            <p style="display: flex;justify-content: left;align-items: center;padding: 8px;font-family: 'Fira Code', monospace;color: #f5f3f1;padding-top: 8px">
                {{ getDiff(data, indextr) }}
            </p>
          </vs-td>
          
          <vs-td :data="data[indextr]" style="background-color: #414141">
<!--            <vs-input class="inputx" v-model.number="timeArr[indextr]" type="number"/>-->
            <p style="display: flex;justify-content: left;align-items: center;padding: 8px;font-family: 'Fira Code', monospace;color: #f5f3f1;padding-top: 8px">
                {{ timeAsHour(getUnused(features, allUnUsedTime, indextr)) }}
            </p>
          </vs-td>

          <vs-td :data="data[indextr]" style="background-color: #414141">
            <p style="display: flex;justify-content: left;align-items: center;padding: 8px;font-family: 'Fira Code', monospace;color: #f5f3f1;padding-top: 8px">
                {{timeAsHour(data[indextr].SpentSum + getUnused(features, allUnUsedTime, indextr)) }}
            </p>
          </vs-td>
        </vs-tr>
      </template>
    </vs-table>
    `
}