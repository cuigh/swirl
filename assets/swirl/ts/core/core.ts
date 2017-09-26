///<reference path="bulma.ts"/>
///<reference path="ajax.ts"/>
///<reference path="validator.ts"/>
///<reference path="form.ts"/>
///<reference path="dispatcher.ts"/>
///<reference path="table.ts"/>
namespace Swirl.Core {
    class BulmaMarker implements ValidationMarker {
        setError($input: JQuery, errors: string[]) {
            let $field = this.getField($input);

            // update input state
            $input.removeClass('is-success').addClass('is-danger');

            // set errors into errors block
            let $errors = $field.find('div.errors');
            if (!$errors.length) {
                $errors = $('<div class="errors"/>').appendTo($field);
            }
            $errors.empty().append($.map(errors, (err: string) => `<p class="help is-danger">${err}</p>`)).show();
        }

        clearError($input: JQuery) {
            let $field = this.getField($input);

            // update input state
            $input.removeClass('is-danger').addClass('is-success');

            // clear errors
            let $errors = $field.find("div.errors");
            $errors.empty().hide();
        }

        reset($input: JQuery): void {
            let $field = this.getField($input);

            // update input state
            $input.removeClass('is-danger is-success');

            // clear errors
            let $errors = $field.find("div.errors");
            $errors.empty().hide();
        }

        private getField($input: JQuery): JQuery {
            let $field = $input.closest(".field");
            if ($field.hasClass("has-addons") || $field.hasClass("is-grouped")) {
                $field = $field.parent();
            }
            return $field;
        }
    }

    /**
     * Initialize bulma adapters
     */
    $(() => {
        // menu burger
        $('.navbar-burger').click(function () {
            let $el = $(this);
            let $target = $('#' + $el.data('target'));
            $el.toggleClass('is-active');
            $target.toggleClass('is-active');
        });

        // Modal
        Modal.initialize();

        // Tab
        Tab.initialize();

        // AJAX
        AjaxRequest.preHandler = opts => opts.trigger && $(opts.trigger).addClass("is-loading");
        AjaxRequest.postHandler = opts => opts.trigger && $(opts.trigger).removeClass("is-loading");

        // Validator
        Validator.setMarker(new BulmaMarker());

        // Form
        Form.automate();
    });
}