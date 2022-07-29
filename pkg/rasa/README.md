## Rasa Local Configuration

### Initial Configuration

#### Activate Conda

```shell
$ source ~/miniforge3/bin/activate
```

#### Create Environment

```shell
(base) $ conda env create -v --name rasaenv -f env.yml
```

### Usage

#### Activate Environment

```shell
(base) $ conda activate rasaenv
```

#### Train the NLU (Natural Language Understanding) Model

```shell
(rasaenv) $ python -m rasa train nlu
```

#### Run the Trained Model

```shell
(rasaenv) $  python -m rasa run --enable-api -m models/{model_name}
```

#### Deactivate Environment

```shell
(rasaenv) $  conda deactivate
```