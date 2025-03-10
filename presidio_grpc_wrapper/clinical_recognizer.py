from presidio_analyzer import EntityRecognizer, RecognizerResult
from transformers import AutoTokenizer, AutoModelForTokenClassification, pipeline

class ClinicalBERTRecognizer(EntityRecognizer):
    def __init__(self):
        # Download the model from Hugging Face's model hub
        model_name = "blaze999/Medical-NER"

        # Load the model and tokenizer
        self.tokenizer = AutoTokenizer.from_pretrained(model_name)
        self.model = AutoModelForTokenClassification.from_pretrained(model_name)

        # Create a pipeline for named entity recognition
        self.ner_pipeline = pipeline("ner", model=self.model, tokenizer=self.tokenizer)

        # Define the supported entities
        self.supported_entities = [
            "BIOLOGICAL_ATTRIBUTE",
            "BIOLOGICAL_STRUCTURE",
            "CLINICAL_EVENT",
            "DISEASE_DISORDER",
            "FAMILY_HISTORY",
            "HISTORY",
            "MEDICATION",
        ]

        super().__init__(supported_entities=self.supported_entities)


    def analyze(self, text, entities, nlp_artifacts=None):
        results = []

        # Perform named entity recognition on the input text
        ner_results = self.ner_pipeline(text)

        for entity in ner_results:
            entity_type = entity["entity"].replace("B-", "").replace("I-", "")

            # Check if the entity type is in  the list of supported entities
            if entity_type in self.supported_entities:
                recognizer_result = RecognizerResult(
                    entity_type=entity_type,
                    start=entity["start"],
                    end=entity["end"],
                    score=entity["score"]
                )

                # Create a RecognizerResult object for the entity
                results.append(recognizer_result)

        return results