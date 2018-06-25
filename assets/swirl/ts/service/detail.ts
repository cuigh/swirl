///<reference path="../core/core.ts" />
namespace Swirl.Service {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;

    export class DetailPage {
        constructor() {
            $("#btn-delete").click(this.deleteService.bind(this));
            $("#btn-scale").click(this.scaleService.bind(this));
            $("#btn-restart").click(this.restartService.bind(this));
            $("#btn-rollback").click(this.rollbackService.bind(this));
        }

        private deleteService(e: JQueryEventObject) {
            let name = $("#h2-name").text().trim();
            Modal.confirm(`Are you sure to remove service: <strong>${name}</strong>?`, "Delete service", (dlg, e) => {
                $ajax.post(`delete`).trigger(e.target).encoder("form").json<AjaxResult>(() => {
                    location.href = "/service/";
                })
            });
        }

        private scaleService(e: JQueryEventObject) {
            let data = {
                count: $("#span-replicas").text().trim(),
            };
            Modal.confirm(`<input name="count" value="${data.count}" class="input" placeholder="Replicas">`, "Scale service", dlg => {
                data.count = dlg.find("input[name=count]").val();
                $ajax.post(`scale`, data).trigger(e.target).encoder("form").json<AjaxResult>(() => {
                    location.reload();
                })
            });
        }

        private rollbackService(e: JQueryEventObject) {
            let name = $("#h2-name").text().trim();
            Modal.confirm(`Are you sure to rollback service: <strong>${name}</strong>?`, "Rollback service", dlg => {
                $ajax.post(`rollback`).trigger(e.target).encoder("form").json<AjaxResult>(() => {
                    location.reload();
                })
            });
        }

        private restartService(e: JQueryEventObject) {
            let name = $("#h2-name").text().trim();
            Modal.confirm(`Are you sure to restart service: <strong>${name}</strong>?`, "Restart service", dlg => {
                $ajax.post(`restart`).trigger(e.target).encoder("form").json<AjaxResult>(() => {
                    location.reload();
                })
            });
        }
    }
}