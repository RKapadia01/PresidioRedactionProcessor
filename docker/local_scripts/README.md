# To Invoke Local Scripts

You need to have a working Presidio environment to run these scripts locally.

```
pip install presidio-analyzer
pip install presidio-anonymizer
python -m spacy download en_core_web_lg
```

## Analyzer Only

```
echo '{"text": "My phone number is 555-123-4567 and my name is Jacob.","language": "en"}' | python ./analyzer.py
```

## Anonymizer Only

```
echo '{"text": "My phone number is 555-123-4567 and my name is Jacob.","analyzer_results": [ ...(from the output of Analyzer) ]}' | python ./anonymizer.py
```

## Combined Analyzer and Anonymizer

```
echo '{"text": "My phone number is 555-123-4567 and my name is Jacob.","language": "en"}' | python ./analyzer_anonymizer.py
```
