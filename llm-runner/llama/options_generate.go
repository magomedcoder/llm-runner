package llama

// WithMaxTokens задает верхнюю границу числа токенов в ответе.
// Помогает ограничить длину генерации и время выполнения запроса.
func WithMaxTokens(n int) GenerateOption {
	return func(c *generateConfig) {
		c.maxTokens = n
	}
}

// WithTemperature задает температуру сэмплирования
// Низкие значения делают ответ более детерминированным, высокие - более вариативным
func WithTemperature(t float32) GenerateOption {
	return func(c *generateConfig) {
		c.temperature = t
	}
}

// WithTopP задает порог nucleus sampling (top-p)
// Модель выбирает токены из "ядра" вероятностного распределения
func WithTopP(p float32) GenerateOption {
	return func(c *generateConfig) {
		c.topP = p
	}
}

// WithTopK ограничивает выбор следующего токена только k наиболее вероятными
// Уменьшение значения делает генерацию более предсказуемой
func WithTopK(k int) GenerateOption {
	return func(c *generateConfig) {
		c.topK = k
	}
}

// WithSeed задает seed генератора случайных чисел
// Одинаковый seed помогает воспроизводить результаты при прочих равных параметрах
func WithSeed(seed int) GenerateOption {
	return func(c *generateConfig) {
		c.seed = seed
	}
}

// WithStopWords задает последовательности, при появлении которых генерация останавливается
// Используется для ранней остановки по маркерам конца ответа
func WithStopWords(words ...string) GenerateOption {
	return func(c *generateConfig) {
		c.stopWords = words
	}
}

// WithDraftTokens задает количество "черновых" токенов для спекулятивной генерации
// Может ускорять вывод при поддержке механизма draft в бэкенде
func WithDraftTokens(n int) GenerateOption {
	return func(c *generateConfig) {
		c.draftTokens = n
	}
}

// WithDebug включает режим отладки генерации
// Обычно добавляет диагностическую информацию в логи или ответ
func WithDebug() GenerateOption {
	return func(c *generateConfig) {
		c.debug = true
	}
}

// WithMinP задает минимальный относительный порог вероятности токенов (min-p)
// Отсекает слишком маловероятные варианты даже при широком top-p/top-k
func WithMinP(p float32) GenerateOption {
	return func(c *generateConfig) {
		c.minP = p
	}
}

// WithTypicalP задает typical sampling (typical-p)
// Смещает выбор к "типичным" токенам, уменьшая слишком редкие отклонения
func WithTypicalP(p float32) GenerateOption {
	return func(c *generateConfig) {
		c.typP = p
	}
}

// WithTopNSigma задает отсечение токенов по отклонению вероятностей (top-n-sigma)
// Ограничивает кандидатов статистически менее правдоподобными вариантами
func WithTopNSigma(sigma float32) GenerateOption {
	return func(c *generateConfig) {
		c.topNSigma = sigma
	}
}

// WithMinKeep задает минимальное число токенов-кандидатов после всех фильтров
// Нужен, чтобы выбор не сужался чрезмерно в сложных распределениях
func WithMinKeep(n int) GenerateOption {
	return func(c *generateConfig) {
		c.minKeep = n
	}
}

// WithRepeatPenalty задает общий штраф за повторение ранее сгенерированных токенов
// Помогает уменьшить зацикливание и однообразие текста
func WithRepeatPenalty(penalty float32) GenerateOption {
	return func(c *generateConfig) {
		c.penaltyRepeat = penalty
	}
}

// WithFrequencyPenalty задает штраф, растущий с частотой появления токена
// Чем чаще токен уже встречался, тем сильнее его вероятность понижается
func WithFrequencyPenalty(penalty float32) GenerateOption {
	return func(c *generateConfig) {
		c.penaltyFreq = penalty
	}
}

// WithPresencePenalty задает штраф за сам факт присутствия токена в истории
// Полезно для повышения тематического разнообразия ответа
func WithPresencePenalty(penalty float32) GenerateOption {
	return func(c *generateConfig) {
		c.penaltyPresent = penalty
	}
}

// WithPenaltyLastN задает размер окна последних токенов для применения штрафов
// Контролирует, на какую часть недавнего контекста смотреть при penalization
func WithPenaltyLastN(n int) GenerateOption {
	return func(c *generateConfig) {
		c.penaltyLastN = n
	}
}

// WithDRYMultiplier задает силу DRY-механизма подавления повторов
// Увеличение значения активнее борется с повторяющимися фрагментами
func WithDRYMultiplier(mult float32) GenerateOption {
	return func(c *generateConfig) {
		c.dryMultiplier = mult
	}
}

// WithDRYBase задает базовый уровень DRY-штрафа
// Работает как опорное значение для расчета итогового подавления повторов
func WithDRYBase(base float32) GenerateOption {
	return func(c *generateConfig) {
		c.dryBase = base
	}
}

// WithDRYAllowedLength задает длину последовательности, которая еще не считается повтором
// Более высокое значение делает DRY-механику мягче
func WithDRYAllowedLength(length int) GenerateOption {
	return func(c *generateConfig) {
		c.dryAllowedLength = length
	}
}

// WithDRYPenaltyLastN задает окно истории токенов для DRY-анализа
// Чем больше окно, тем дальше в прошлое модель ищет повторы
func WithDRYPenaltyLastN(n int) GenerateOption {
	return func(c *generateConfig) {
		c.dryPenaltyLastN = n
	}
}

// WithDRYSequenceBreakers задает токены/строки, которые разрывают DRY-последовательность
// Позволяет не считать повторами фрагменты, разделенные особыми маркерами
func WithDRYSequenceBreakers(breakers ...string) GenerateOption {
	return func(c *generateConfig) {
		c.drySequenceBreakers = breakers
	}
}

// WithDynamicTemperature задает диапазон и экспоненту динамической температуры
// Температура адаптируется по ходу генерации, изменяя баланс креативности и стабильности
func WithDynamicTemperature(tempRange, exponent float32) GenerateOption {
	return func(c *generateConfig) {
		c.dynatempRange = tempRange
		c.dynatempExponent = exponent
	}
}

// WithXTC задает параметры XTC-сэмплирования: вероятность и порог
// Используется для дополнительной фильтрации/коррекции выбора токенов
func WithXTC(probability, threshold float32) GenerateOption {
	return func(c *generateConfig) {
		c.xtcProbability = probability
		c.xtcThreshold = threshold
	}
}

// WithMirostat включает и задает версию алгоритма Mirostat
// Алгоритм поддерживает целевой уровень неожиданности (энтропии) в ответе
func WithMirostat(version int) GenerateOption {
	return func(c *generateConfig) {
		c.mirostat = version
	}
}

// WithMirostatTau задает целевой уровень неожиданности (tau) для Mirostat
// Влияет на "смелость" выбора токенов при активном Mirostat
func WithMirostatTau(tau float32) GenerateOption {
	return func(c *generateConfig) {
		c.mirostatTau = tau
	}
}

// WithMirostatEta задает скорость адаптации (eta) в алгоритме Mirostat
// Более высокие значения быстрее корректируют распределение токенов
func WithMirostatEta(eta float32) GenerateOption {
	return func(c *generateConfig) {
		c.mirostatEta = eta
	}
}

// WithNPrev задает число предыдущих токенов, учитываемых в некоторых сэмплинг-правилах
// Влияет на то, насколько "далеко назад" смотрят механизмы штрафов и эвристик
func WithNPrev(n int) GenerateOption {
	return func(c *generateConfig) {
		c.nPrev = n
	}
}

// WithNProbs задает количество top-вероятностей токенов в выходных данных
// Полезно для отладки, анализа и отображения альтернативных продолжений
func WithNProbs(n int) GenerateOption {
	return func(c *generateConfig) {
		c.nProbs = n
	}
}

// WithIgnoreEOS включает или отключает игнорирование EOS-токена
// При включении модель может продолжать генерацию даже после потенциального конца
func WithIgnoreEOS(ignore bool) GenerateOption {
	return func(c *generateConfig) {
		c.ignoreEOS = ignore
	}
}
