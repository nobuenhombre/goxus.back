package settingsdomain

const (
	SettingsTypeInputTextField        = "inputTextField"
	SettingsTypeInputPasswordField    = "inputPasswordField"
	SettingsTypeInputIntNumberField   = "inputIntNumberField"
	SettingsTypeInputFloatNumberField = "inputFloatNumberField"
	SettingsTypeTextareaField         = "textareaField"
	SettingsTypeInputIntSlider        = "inputIntSlider"
	SettingsTypeInputIntSliderRange   = "inputIntSliderRange"
	SettingsTypeSwitch                = "switch"
	SettingsTypeListChecks            = "listChecks"
	SettingsTypeListRadios            = "listRadios"
	SettingsTypeSelectSimple          = "selectSimple"
	SettingsTypeSelectWithSearch      = "selectWithSearch"
)

// SettingsTypeInputTextField = "inputTextField"
// example saved json value
// {
//    "value": "hello world"
// }
// example saved json available_values: null
// example saved json default_value:
// {
//    "value": "hello world"
// }

//------------------------------------------------------
// SettingsTypeInputPasswordField = "inputPasswordField"
// example saved json value
// {
//    "value": "my password"
// }
// example saved json available_values: null
// example saved json default_value:
// {
//    "value": "123"
// }

//------------------------------------------------------
// SettingsTypeInputIntNumberField = "inputIntNumberField"
// example saved json value
// {
//    "value": 5
// }
// example saved json available_values: null
// example saved json default_value:
// {
//    "value": 0
// }

//------------------------------------------------------
// SettingsTypeInputFloatNumberField = "inputFloatNumberField"
// example saved json value
// {
//    "value": 5.5
// }
// example saved json available_values: null
// example saved json default_value:
// {
//    "value": 0.0
// }

//------------------------------------------------------
// SettingsTypeTextareaField = "textareaField"
// example saved json value
// {
//    "value": "hello world"
// }
// example saved json available_values: null
// example saved json default_value:
// {
//    "value": "hello world"
// }

//------------------------------------------------------
// SettingsTypeInputIntSlider = "inputIntSlider"
// example saved json value
// {
//    "value": 5
// }
// example saved json available_values: null
// example saved json default_value:
// {
//    "value": 0
// }

//------------------------------------------------------
// SettingsTypeInputIntSliderRange = "inputIntSliderRange"
// example saved json value
// {
//    "value": {
//      "start": 3,
//      "end": 5
//   }
// }
// example saved json available_values:
// {
//    "value": {
//      "min": 0,
//      "max": 10
//   }
// }
// example saved json default_value:
// {
//    "value": {
//      "start": 0,
//      "end": 10
//   }
// }

//------------------------------------------------------
// SettingsTypeSwitch = "switch"
// {
//    "value": true
// }
// example saved json available_values: null
// example saved json default_value:
// {
//    "value": false
// }

//------------------------------------------------------
// SettingsTypeListChecks = "listChecks"
// {
//    "value": {
//      "1": true,
//      "2": false,
//      "3": true,
//   }
// }
// example saved json available_values:
// {
//    "value": {
//      "1": "Hello world my check A",
//      "2": "Hello world my check B",
//      "3": "Hello world my check C",
//   }
// }
// example saved json default_value:
// {
//    "value": {
//      "1": false,
//      "2": false,
//      "3": false,
//   }
// }

//------------------------------------------------------
// SettingsTypeListRadios = "listRadios"
// {
//    "value": 1
// }
// example saved json available_values:
// {
//    "value": {
//      "1": "Hello world my radio A",
//      "2": "Hello world my radio B",
//      "3": "Hello world my radio C",
//   }
// }
// example saved json default_value:
// {
//    "value": 1
// }

//------------------------------------------------------
// SettingsTypeSelectSimple = "selectSimple"
// {
//    "value": 3
// }
// example saved json available_values:
// {
//    "value": {
//      "1": "Hello world my select item A",
//      "2": "Hello world my select item B",
//      "3": "Hello world my select item C",
//   }
// }
// example saved json default_value:
// {
//    "value": 3
// }

//------------------------------------------------------
// SettingsTypeSelectWithSearch = "selectWithSearch"
// {
//    "value": 3
// }
// example saved json available_values:
// {
//    "value": {
//      "1": "Hello world my select item A",
//      "2": "Hello world my select item B",
//      "3": "Hello world my select item C",
//   }
// }
// example saved json default_value:
// {
//    "value": 3
// }
