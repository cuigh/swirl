///<reference path="../core/core.ts" />
namespace Swirl.Service {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Table = Swirl.Core.ListTable;

    export class ListPage {
        private table: Table;

        constructor() {
            this.table = new Table("#table-items");

            // bind events
            this.table.on("delete-service", this.deleteService.bind(this))
                .on("scale-service", this.scaleService.bind(this))
                .on("rollback-service", this.rollbackService.bind(this))
                .on("restart-service", this.restartService.bind(this));
        }

        private deleteService(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let name = $tr.find("td:eq(0)").text().trim();
            Modal.confirm(`Are you sure to remove service: <strong>${name}</strong>?`, "Delete service", (dlg, e) => {
                $ajax.post(`${name}/delete`).trigger(e.target).encoder("form").json<AjaxResult>(() => {
                    $tr.remove();
                    dlg.close();
                })
            });
        }

        private scaleService(e: JQueryEventObject) {
            let $target = $(e.target),
                $tr = $target.closest("tr"),
                name = $tr.find("td:eq(0)").text().trim();
            let data = {
                count: $target.data("replicas"),
            };
            Modal.confirm(`<input name="count" value="${data.count}" class="input" placeholder="Replicas">`, "Scale service", dlg => {
                data.count = dlg.find("input[name=count]").val();
                $ajax.post(`${name}/scale`, data).encoder("form").json<AjaxResult>(() => {
                    location.reload();
                })
            });
        }

        private rollbackService(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr"),
                name = $tr.find("td:eq(0)").text().trim();
            Modal.confirm(`Are you sure to rollback service: <strong>${name}</strong>?`, "Rollback service", dlg => {
                $ajax.post(`${name}/rollback`).encoder("form").json<AjaxResult>(() => {
                    dlg.close();
                })
            });
        }

        private restartService(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr"),
                name = $tr.find("td:eq(0)").text().trim();
            Modal.confirm(`Are you sure to restart service: <strong>${name}</strong>?`, "Restart service", dlg => {
                $ajax.post(`${name}/restart`).encoder("form").json<AjaxResult>(() => {
                    dlg.close();
                })
            });
        }
    }
}