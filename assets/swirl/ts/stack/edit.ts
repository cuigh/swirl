///<reference path="../core/core.ts" />
namespace Swirl.Stack {
    import Validator = Swirl.Core.Validator;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Notification = Swirl.Core.Notification;
    import ValidationRule = Swirl.Core.ValidationRule;

    class ContentRequiredRule implements ValidationRule {
        validate($form: JQuery, $input: JQuery, arg?: string): {ok: boolean, error?: string} {
            let el = <HTMLInputElement>$input[0];
            if ($("#type-" + arg).prop("checked")) {
                console.log(el.value);
                return {ok: el.checkValidity ? el.checkValidity() : true, error: el.validationMessage};
            }
            return {ok: true}
        }
    }

    export class EditPage {
        private editor: any;

        constructor() {
            Validator.register("content", new ContentRequiredRule(), "");

            this.editor = CodeMirror.fromTextArea($("#txt-content")[0], {lineNumbers: true});

            $("#file-content").change(e => {
                let file = <HTMLInputElement>e.target;
                if (file.files.length > 0) {
                    $('#filename').text(file.files[0].name);
                }
            });
            $("#type-input,#type-upload").click(e => {
                let type = $(e.target).val();
                $("#div-input").toggle(type == "input");
                $("#div-upload").toggle(type == "upload");
            });
            $("#btn-submit").click(this.submit.bind(this))
        }

        private submit(e: JQueryEventObject) {
            this.editor.save();

            let results = Validator.bind("#div-form").validate();
            if (results.length > 0) {
                return;
            }

            let data = new FormData();
            data.append('name', $("#name").val());
            if ($("#type-input").prop("checked")) {
                data.append('content', $('#txt-content').val());
            } else {
                let file = <HTMLInputElement>$('#file-content')[0];
                data.append('content', file.files[0]);
            }

            let url = $(e.target).data("url") || "";
            $ajax.post(url, data).encoder("none").trigger(e.target).json((r: AjaxResult) => {
                if (r.success) {
                    location.href = "/stack/"
                } else {
                    Notification.show("danger", `FAILED: ${r.message}`);
                }
            })
        }
    }
}

declare var CodeMirror: any;