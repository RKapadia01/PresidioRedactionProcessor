import sys
import json
from presidio_analyzer import AnalyzerEngine
from presidio_anonymizer import AnonymizerEngine

def main():
    input_data = sys.stdin.read()

    try:
        analyzer_params = json.loads(input_data)
        analyzer = AnalyzerEngine()
        analyzer_results = analyzer.analyze(**analyzer_params)

        anonymizer = AnonymizerEngine()
        results = anonymizer.anonymize(text=analyzer_params["text"], analyzer_results=analyzer_results)

        print(results.text)
    except json.JSONDecodeError as e:
        print(f"Error: Invalid JSON input. {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()