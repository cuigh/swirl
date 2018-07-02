///<reference path="../core/core.ts" />
namespace Swirl.Metric {
    import EditTable = Swirl.Core.EditTable;

    class MetricTable extends EditTable {
        protected render(): string {
            return `<tr>
                <td>
                  <input name="metrics[${this.index}].legend" class="input is-small" placeholder="Legend expression for dataset, e.g. ${name}">
                </td>
                <td>
                  <input name="metrics[${this.index}].query" class="input is-small" placeholder="Prometheus query expression, for service dashboard, you can use '$\{service\}' variable">
                </td>
                <td>
                  <a class="button is-small is-outlined is-danger" data-action="delete-metric">
                    <span class="icon is-small">
                      <i class="far fa-trash-alt"></i>
                    </span>
                  </a>
                </td>
              </tr>`;
        }
    }

    export class EditPage {
        constructor() {
            new MetricTable("#table-metrics");
        }
    }
}