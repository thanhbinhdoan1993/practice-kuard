import React from 'react';
import Form from 'react-jsonschema-form';
import fetchError from './fetcherror';

const schema = {
    "$schema": "http://json-schema.org/draft-04/schema#",
    "type": "object",
    "properties": {
        "enable": {
            "title": "Enabled?",
            "type": "boolean"
        },
        "exitOnComplete": {
            "title": "Exit server on completion?",
            "type": "boolean"
        },
        "exitCode": {
            "title": "Exit code when exiting. 0 is success.",
            "type": "integer"
        },
        "numToGen": {
            "title": "Number of keys to generate. 0 is infinite.",
            "type": "integer"
        },
        "timeToRun": {
            "title": "Time to run, in seconds. 0 is infinite.",
            "type": "integer"
        },
        "memQServer": {
            "title": "Base URL of the MemQ server to draw from. Can be http://localhost:8080/memq/server.",
            "type": "string"
        },
        "memQQueue": {
            "title": "The Queue to pull work items from.",
            "type": "string"
        }
    }
};