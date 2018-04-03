///<reference path="core/core.ts" />
///<reference path="core/chart.ts" />
namespace Swirl {
    import ChartDashboard = Swirl.Core.ChartDashboard;

    export class IndexPage {
        private dashboard: ChartDashboard;

        constructor() {
            this.dashboard = new ChartDashboard("#div-charts", window.charts, {name: "home"});
            dragula([$('#div-charts').get(0)], {
                moves: function (el, container, handle): boolean {
                    return $(handle).closest('a.drag').length > 0;
                    // return handle.classList.contains('drag');
                }
            });
        }
    }
}